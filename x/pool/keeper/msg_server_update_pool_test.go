package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govV1Types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// Pool
	"github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_update_pool.go

* Invalid authority (transaction)
* Invalid authority (proposal)
* Update first pool
* Update first pool partially
* Update another pool
* Update pool with invalid json payload

*/

var _ = Describe("msg_server_update_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	gov := s.App().GovKeeper.GetGovernanceAccount(s.Ctx()).GetAddress().String()
	votingPeriod := s.App().GovKeeper.GetVotingParams(s.Ctx()).VotingPeriod

	BeforeEach(func() {
		s = i.NewCleanChain()

		s.App().PoolKeeper.AppendPool(s.Ctx(), types.Pool{
			Name:                     "",
			Runtime:                  "",
			Logo:                     "",
			Config:                   "",
			UploadInterval:           0,
			OperatingCost:            0,
			MinDelegation:            0,
			MaxBundleSize:            0,
			CurrentStorageProviderId: 0,
			CurrentCompressionId:     0,
			Protocol:                 &types.Protocol{},
			UpgradePlan:              &types.UpgradePlan{},
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Invalid authority (transaction)", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: i.DUMMY[0],
			Id:        0,
			Payload:   "{\"Name\":\"TestPool\",\"Runtime\":\"@kyve/test\",\"Logo\":\"ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU\",\"Config\":\"ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0\",\"StartKey\":\"0\",\"UploadInterval\":60,\"OperatingCost\":10000,\"MinDelegation\":\"100000000000\",\"MaxBundleSize\":100,\"Version\":\"0.0.0\",\"Binaries\":\"{}\",\"StorageProviderId\":2,\"CompressionId\":1}",
		}

		// ACT
		_, err := s.RunTx(msg)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Invalid authority (proposal)", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: i.DUMMY[0],
			Id:        0,
			Payload:   "{\"Name\":\"TestPool\",\"Runtime\":\"@kyve/test\",\"Logo\":\"ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU\",\"Config\":\"ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0\",\"StartKey\":\"0\",\"UploadInterval\":60,\"OperatingCost\":10000,\"MinDelegation\":\"100000000000\",\"MaxBundleSize\":100,\"Version\":\"0.0.0\",\"Binaries\":\"{}\",\"StorageProviderId\":2,\"CompressionId\":1}",
		}

		proposal, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, err := s.RunTx(&proposal)

		// ASSERT
		Expect(err).To(HaveOccurred())
	})

	It("Update first pool", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        0,
			Payload:   "{\"Name\":\"TestPool\",\"Runtime\":\"@kyve/test\",\"Logo\":\"ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU\",\"Config\":\"ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0\",\"StartKey\":\"0\",\"UploadInterval\":60,\"OperatingCost\":10000,\"MinDelegation\":100000000000,\"MaxBundleSize\":100,\"Version\":\"0.0.0\",\"Binaries\":\"{}\",\"StorageProviderId\":2,\"CompressionId\":1}",
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool).To(Equal(types.Pool{
			Id:             0,
			Name:           "TestPool",
			Runtime:        "@kyve/test",
			Logo:           "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:         "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:       "",
			CurrentKey:     "",
			CurrentSummary: "",
			CurrentIndex:   0,
			TotalBundles:   0,
			UploadInterval: 60,
			OperatingCost:  10000,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Disabled:       false,
			Funders:        nil,
			TotalFunds:     0,
			Protocol: &types.Protocol{
				Version:     "",
				Binaries:    "",
				LastUpgrade: 0,
			},
			UpgradePlan: &types.UpgradePlan{
				Version:     "",
				Binaries:    "",
				ScheduledAt: 0,
				Duration:    0,
			},
			CurrentStorageProviderId: 2,
			CurrentCompressionId:     1,
		}))
	})

	It("Update first pool partially", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        0,
			Payload:   "{\"Name\":\"TestPool\",\"Runtime\":\"@kyve/test\"}",
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(pool).To(Equal(types.Pool{
			Id:             0,
			Name:           "TestPool",
			Runtime:        "@kyve/test",
			Logo:           "",
			Config:         "",
			StartKey:       "",
			CurrentKey:     "",
			CurrentSummary: "",
			CurrentIndex:   0,
			TotalBundles:   0,
			UploadInterval: 0,
			OperatingCost:  0,
			MinDelegation:  0,
			MaxBundleSize:  0,
			Disabled:       false,
			Funders:        nil,
			TotalFunds:     0,
			Protocol: &types.Protocol{
				Version:     "",
				Binaries:    "",
				LastUpgrade: 0,
			},
			UpgradePlan: &types.UpgradePlan{
				Version:     "",
				Binaries:    "",
				ScheduledAt: 0,
				Duration:    0,
			},
			CurrentStorageProviderId: 0,
			CurrentCompressionId:     0,
		}))
	})

	It("Update another pool", func() {
		// ARRANGE
		s.App().PoolKeeper.AppendPool(s.Ctx(), types.Pool{
			Name:                     "",
			Runtime:                  "",
			Logo:                     "",
			Config:                   "",
			UploadInterval:           0,
			OperatingCost:            0,
			MinDelegation:            0,
			MaxBundleSize:            0,
			CurrentStorageProviderId: 0,
			CurrentCompressionId:     0,
			Protocol:                 &types.Protocol{},
			UpgradePlan:              &types.UpgradePlan{},
		})

		// ACT
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        1,
			Payload:   "{\"Name\":\"TestPool2\",\"Runtime\":\"@kyve/test\",\"Logo\":\"ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU\",\"Config\":\"ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0\",\"StartKey\":\"0\",\"UploadInterval\":60,\"OperatingCost\":10000,\"MinDelegation\":100000000000,\"MaxBundleSize\":100,\"Version\":\"0.0.0\",\"Binaries\":\"{}\",\"StorageProviderId\":2,\"CompressionId\":1}",
		}

		p, v := BuildGovernanceTxs(s, []sdk.Msg{msg})

		_, submitErr := s.RunTx(&p)
		_, voteErr := s.RunTx(&v)

		s.CommitAfter(*votingPeriod)
		s.Commit()

		// ASSERT
		proposal, _ := s.App().GovKeeper.GetProposal(s.Ctx(), 1)

		Expect(submitErr).To(Not(HaveOccurred()))
		Expect(voteErr).To(Not(HaveOccurred()))

		Expect(proposal.Status).To(Equal(govV1Types.StatusPassed))

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 1)
		Expect(pool).To(Equal(types.Pool{
			Id:             1,
			Name:           "TestPool2",
			Runtime:        "@kyve/test",
			Logo:           "ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU",
			Config:         "ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0",
			StartKey:       "",
			CurrentKey:     "",
			CurrentSummary: "",
			CurrentIndex:   0,
			TotalBundles:   0,
			UploadInterval: 60,
			OperatingCost:  10000,
			MinDelegation:  100 * i.KYVE,
			MaxBundleSize:  100,
			Disabled:       false,
			Funders:        nil,
			TotalFunds:     0,
			Protocol: &types.Protocol{
				Version:     "",
				Binaries:    "",
				LastUpgrade: 0,
			},
			UpgradePlan: &types.UpgradePlan{
				Version:     "",
				Binaries:    "",
				ScheduledAt: 0,
				Duration:    0,
			},
			CurrentStorageProviderId: 2,
			CurrentCompressionId:     1,
		}))
	})

	It("Update pool with invalid json payload", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        1,
			Payload:   "invalid_json_payload\",\"Runtime\":\"@kyve/test\",\"Logo\":\"ar://Tewyv2P5VEG8EJ6AUQORdqNTectY9hlOrWPK8wwo-aU\",\"Config\":\"ar://DgdB-2hLrxjhyEEbCML__dgZN5_uS7T6Z5XDkaFh3P0\",\"StartKey\":\"0\",\"UploadInterval\":60,\"OperatingCost\":10000,\"MinDelegation\":100000000000,\"MaxBundleSize\":100,\"Version\":\"0.0.0\",\"Binaries\":\"{}\",\"StorageProviderId\":2,\"CompressionId\":1}",
		}

		p, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_ = s.RunTxError(&p)
		s.Commit()

		// ASSERT
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(found).To(BeTrue())
		Expect(pool.Name).To(BeEmpty())
	})

	It("Update pool with invalid UploadInterval", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        1,
			Payload:   "{\"UploadInterval\": 0}",
		}

		p, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_ = s.RunTxError(&p)
		s.Commit()

		// ASSERT
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(found).To(BeTrue())
		Expect(pool.Name).To(BeEmpty())
	})

	It("Update pool with invalid UploadInterval", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        1,
			Payload:   "{\"UploadInterval\": 0}",
		}

		p, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_ = s.RunTxError(&p)
		s.Commit()

		// ASSERT
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(found).To(BeTrue())
		Expect(pool.Name).To(BeEmpty())
	})

	It("Update pool with invalid OperatingCost", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        1,
			Payload:   "{\"OperatingCost\": 0}",
		}

		p, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_ = s.RunTxError(&p)
		s.Commit()

		// ASSERT
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(found).To(BeTrue())
		Expect(pool.Name).To(BeEmpty())
	})

	It("Update pool with invalid MinDelegation", func() {
		// ARRANGE
		msg := &types.MsgUpdatePool{
			Authority: gov,
			Id:        1,
			Payload:   "{\"MinDelegation\": 0}",
		}

		p, _ := BuildGovernanceTxs(s, []sdk.Msg{msg})

		// ACT
		_ = s.RunTxError(&p)
		s.Commit()

		// ASSERT
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(found).To(BeTrue())
		Expect(pool.Name).To(BeEmpty())
	})
})
