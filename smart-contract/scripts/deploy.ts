import { ethers } from "hardhat";

async function main() {
    const [deployer] = await ethers.getSigners();

    const Factory = await ethers.getContractFactory("L8TrafficPayment");
    const receiver = process.env.RECEIVER_ADDRESS || deployer.address;
    const contract = await Factory.deploy(receiver);

    await contract.waitForDeployment();

    console.log("Contract:", await contract.getAddress());
    console.log("Owner:", contract.owner);
    console.log("Receiver:", contract.receiver);
}

main();
