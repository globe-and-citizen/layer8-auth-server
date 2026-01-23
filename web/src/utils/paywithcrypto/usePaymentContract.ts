import { useWriteContract, useConnection } from '@wagmi/vue'
import { parseEther } from 'viem'
import artifact from "../../../../smart-contract/abi/L8TrafficPayment.json"
import { config } from "@/config";

const CONTRACT_ADDRESS = config.CONTRACT_ADDRESS as `0x${string}`


export function useTrafficContract() {
  const connection = useConnection()

  const {
    writeContract,
    data: hash,
    isPending,
    isSuccess,
    error: writeError
  } = useWriteContract()

  /**
   * @param clientID - The string ID from your solidity function
   * @param amountInEth - The amount of SepoliaETH to send (e.g., "0.01")
   */
  async function payTraffic(clientID: string, amountInEth: string) {
    writeContract({
      address: CONTRACT_ADDRESS,
      abi: artifact.abi,
      functionName: 'pay', // Matches your Solidity function
      args: [clientID],    // Matches calldata clientID
      value: parseEther(amountInEth), // Converts "0.01" ETH to Wei automatically
    })
  }

  return {
    connection,
    payTraffic,
    isPending,
    isSuccess,
    writeError,
    hash
  }
}
