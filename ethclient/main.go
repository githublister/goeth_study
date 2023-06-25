package main

import (
	//"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"sync"

	//"github.com/ethereum/go-ethereum/core/types"

	//"github.com/ethereum/go-ethereum/crypto"
	"math"

	//"github.com/holiman/uint256"

	//"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/common/hexutil"

	//"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
	"strings"
	"time"
)

//根据  小数位数确认 和分割显示
func formatBigIntWithDecimals(n *big.Int, decimals uint8) string {
	divisor := new(big.Float).SetFloat64(math.Pow10(int(decimals)))
	nFloat := new(big.Float).SetInt(n)
	result := new(big.Float).Quo(nFloat, divisor)

	// 使用 'f' 格式，没有指数，使用 -1 以动态调整精度（保证不会丢失精度）
	resultStr := result.Text('f', -1)

	parts := strings.Split(resultStr, ".")
	wholePart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}

	var out strings.Builder

	l := len(wholePart)
	for i, r := range wholePart {
		out.WriteRune(r)
		if (l-i-1)%3 == 0 && i != l-1 {
			out.WriteRune(',')
		}
	}

	out.WriteString(decimalPart)

	return out.String()
}


//它将区块的Unix时间戳转换为UTC+8（亚洲/上海）时区的时间：
func convertBlockTimeToUTC8(blockTime uint64) time.Time {
	// 从区块的Unix时间戳生成时间
	t := time.Unix(int64(blockTime), 0)

	// 加载UTC+8时区
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// 转换时间到UTC+8时区
	t = t.In(loc)

	// 返回转换后的时间
	return t
}


//通过交易ash 获取 block hash
func getBlockHashFromTxHash(txHash common.Hash) (common.Hash, error) {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return common.Hash{}, fmt.Errorf("Failed to retrieve transaction: %v", err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return common.Hash{}, fmt.Errorf("Failed to retrieve transaction receipt: %v", err)
	}

	return receipt.BlockHash, nil
}

/**
判断地址是否为合约 地址
 */
func checkIfContractAddress(client *ethclient.Client, address common.Address) bool {
	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return false
	}

	return len(bytecode) > 0
}


// 使用您提供的以太坊节点地址
//const infuraURL = "https://eth-mainnet.g.alchemy.com/v2/ovZAGPQlWJlYIoXzD1ifi4HDppTT4QS9"
//const infuraURL = "https://mainnet.infura.io/v3/a583f941f78049a3ab3ebc43790d8825"
const infuraURL = "https://eth-mainnet-public.unifra.io"

//eth测试网节点
const ethGoerliRpc = "https://eth-goerli.g.alchemy.com/v2/qHgbnaritObcCPa5ETXrjQnDMQbIRkAO"

//abi 测试code
const contractABI = `[{"inputs":[{"internalType":"uint256","name":"_totalSupply","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_address","type":"address"},{"internalType":"bool","name":"_isBlacklisting","type":"bool"}],"name":"blacklist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"blacklists","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"value","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"limited","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxHoldingAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"minHoldingAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bool","name":"_limited","type":"bool"},{"internalType":"address","name":"_uniswapV2Pair","type":"address"},{"internalType":"uint256","name":"_maxHoldingAmount","type":"uint256"},{"internalType":"uint256","name":"_minHoldingAmount","type":"uint256"}],"name":"setRule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"uniswapV2Pair","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`
//测试合约地址
const contractAddress = "0x6982508145454ce325ddbe47a25d4ec3d2311933" // Example contract address, replace with your own

//用户的钱包地址
const oeaAddress ="0x9d491a39111a2bcec8408a889ca601e123bc770a"

//为了方便后面的测试 把client封装成一个公用的方法

type Client struct {
	client *ethclient.Client
}

// NewClient function
func NewClient() *Client {
	client, err := ethclient.Dial(infuraURL )
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	return &Client{client: client}
}

func main() {

	/*
	`Dial`是最基本的用法，用于直接连接到以太坊网络。当你不需要关心连接过程中的超时或取消，只是简单地需要一个客户端来进行以太坊交互时，这是最简单的方法。
	 */
	//callDial()

	/*
	如果你的应用需要更精细的控制连接过程，比如设置连接超时或在某个阶段取消连接，那么DialContext会很有用。例如，你的应用在启动时连接到以太坊网络，但如果连接过程过于漫长，你可能希望在一段时间后超时并取消连接，这样你的应用可以快速失败并尝试其他操作
	 */
	//callDialContext()

	/*
	如果你已经有一个rpc.Client实例，并希望在其基础上创建以太坊客户端，那么可以使用NewClient。一个场景是，如果你的应用既需要使用以太坊特性（通过ethclient.Client），又需要执行某些原始的JSON-RPC调用（通过rpc.Client），那么你可能会先创建一个rpc.Client，然后再用它创建一个ethclient.Client。
	 */
	//callNewClient()

	/*
	在以太坊中，BalanceAt方法是ethclient包中Client结构体的一个方法，它允许你查询一个特定账户在某个区块高度的ETH余额。
	 */
	//callBalanceAt()

	/*
	BlockByHash是ethclient.Client提供的一个方法，它用于从以太坊区块链中获取一个特定的区块。它接收两个参数：一个context.Context对象，以及你想要查询的区块的哈希值。它会返回一个types.Block类型的对象。
	 */
	//callBlockByHash()

	/**
	BlockByNumber 是一个用于获取特定区块的函数，通过传入一个包含区块号的 big.Int 对象作为参数。如果你想要获取最新的区块，你可以传入 nil 作为参数。
	 */
	//callBlockByNumber()

	/*
	BlockNumber是ethclient.Client的一个方法，它返回最新的区块编号。你可以使用这个方法来获取以太坊网络当前的区块高度。
	 */
	//callBlockNumber()

	/**
	CallContract 是一个功能强大的方法，它允许你从一个特定的以太坊地址调用一个智能合约的方法，并且返回该方法的返回值。CallContract 对于只读操作非常有用，因为它不需要任何 gas 或者交易确认。
	 */
	//callCallContract()

	/**
	使用 CallContractAtHash 函数来获取特定区块（通过区块hash）上智能合约的状态
	 */
	//callCallContractAtHash()
	/**
	方法是以太坊客户端库中 Client 结构体提供的一个函数。它用于获取当前链的链ID，用于事务回放保护。
	 */
	//callChainID()

	/**
	Client 方法是以太坊客户端库中 Client 结构体提供的一个函数。它用于获取底层的 RPC 客户端对象。
	 */
	//callClient()


	/**
	CodeAt 是以太坊客户端库中 Client 结构体提供的一个方法，用于获取智能合约的字节码（bytecode）。
	 */
	//callCodeAt()


	/**
	EstimateGas 方法是以太坊智能合约中的一个函数，用于估算执行特定交易或调用合约所需的燃料（gas）消耗量。燃料消耗量是以太坊中衡量计算和交易成本的单位
	 */
	//callEstimateGas()

	/**
	FeeHistory 是一个用于查询燃料费用历史的概念。它提供了过去一段时间内燃料费用的统计信息，以帮助用户了解当前以太坊网络上的平均燃料费用情况。
	 */

	//callFeeHistory()

	/*
	FilterLogs 是 Ethereum Go 客户端（go-ethereum）中用于过滤并获取特定合约事件日志的方法。在智能合约中，当特定操作（例如转账或状态更新）发生时，会发出一个事件。这些事件在区块链上公开记录，任何人都可以查询。
	 */
	//callFilterLogs()

	/**
	，它用于获取指定哈希值对应的区块头。这个函数接受一个上下文对象和哈希值作为参数，返回与哈希值对应的区块头和可能的错误。
	 */
	//callHeaderByHash()

	//callHeaderByNumber()

	//callNetworkID()

	/*
	NonceAt方法接收一个Ethereum地址和一个可选的区块号，并返回该地址在指定区块中的nonce。如果区块号是nil，则该方法返回最新区块的nonce。
	 */
	//callNonceAt()

	/**

	PeerCount是go-ethereum库（ethclient）提供的一个方法，用于获取当前Ethereum节点连接的其他节点数量。

	在分布式网络如Ethereum中，节点之间的连接关系很重要。这些连接（或称为对等体）使得节点可以传播交易和区块，使得整个网络能够正常运行。如果一个节点没有对等体，那么它将不能接收到新的交易和区块。
	 */
	//callPeerCount()

	/**
	PendingBalanceAt 方法是一个非常实用的功能，它可以帮助你获取一个特定地址在当前挂起状态的余额。

	该方法的主要用途是获取给定地址在当前还未被挖矿所确认的交易中的余额。这个方法特别有用在高网络负载的时候，那时候可能有许多交易还在等待被确认。
	 */
	//callPendingBalanceAt()

	/**
	执行一个调用（读取）合约的操作，不会产生任何链上的交易。该调用是在挂起的状态下执行，即考虑了当前挂起的交易。
	 */
	//callPendingCallContract()


	/**
	PendingCodeAt：获取给定地址在挂起状态下的代码。这在跟踪合约部署过程中很有用
	 */
	//callPendingCodeAt()

	/**
	PendingNonceAt这个函数的主要作用就是帮助你获取你的账户的下一个有效的nonce值。在手动创建和发送以太坊交易时
	 */
	//callPendingNonceAt()

	/**
	在以太坊中，每个账户有关联的存储空间，可以保存特定的状态信息。PendingStorageAt函数就是用来从某个特定的合约的存储空间获取数据。
	*/
	//callPendingStorageAt()


	/**
	PendingTransactionCount这个方法返回当前节点的事务池中正在等待的事务数量，也就是整个网络中的未确认事务，
	 */
	//callPendingTransactionCount()

	/**
	 函数是用于将一个已签名的交易发送到以太坊网络
	从go-ethereum v1.10.0 开始，types.NewTransaction 的确已经被标记为废弃，并且被 types.NewTx 和 types.NewLegacyTx 替代。
	 */
	//callSendTransaction()

	/**
	StorageAt方法查询的是已经被挖矿并且添加到区块链中的特定区块的状态。这是一个已经发生并且被网络达成共识的状态。
	 */
	//callStorageAt()

	/**
	用于订阅区块链中符合指定过滤条件的日志事件。
	 */
	//callSubscribeFilterLogs()


	/**
	，用于订阅以太坊区块链的新区块头（block header）更新事件。当有新的区块被挖出并添加到区块链上时，订阅者将接收到相应的通知。
	 */
	//callSubscribeNewHead()


	/**

	用于获取当前建议的燃气价格（gas price）。燃气价格是以太坊交易的费用
	 */
	//callSuggestGasPrice()

	/**
	，用于获取当前建议的燃气小费上限。
	 */
	//callSuggestGasTipCap()

	/**
	用于获取节点的同步进度。该方法返回一个包含当前同步状态的结构体，其中包括已下载的区块数、目标区块数、已知的最新区块数等信息。
	 */
	//callcallSyncProgress()


	/**
	，TransactionByHash 方法用于通过交易哈希返回对应的交易信息。
	 */
	//callTransactionByHash()

	/**
	TransactionCount 方法用于获取指定区块的交易数量。它接受一个上下文对象和一个区块哈希作为参数，并返回该区块中的交易数量。
	 */
	//callTransactionCount()

	/**
	TransactionInBlock方法提供了一种从给定块中获取单个事务的便捷方式，帮助您在以太坊区块链上进行事务验证、查询和分析。
	 */
	//callTransactionInBlock()

	/**
	它包含了一个交易的执行结果和相关的信息。通过获取交易的收据，您可以获得以下信息：
	在使用 TransactionReceipt 方法时，您需要提供一个交易的哈希值作为参数，以获取相应交易的收据信息。
	 */
	//callTransactionReceipt()

	/**
	TransactionSender  根据交易hash 获取交易的数据
	 */
	callTransactionSender()
}

func callDial()  {

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fmt.Println("We have a connection")
	_ = client // use the client
}

func callDialContext()  {

	ctx,cancel :=context.WithTimeout(context.Background(),10)
	defer  cancel()
	client, err := ethclient.DialContext(ctx, infuraURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("We have a dialcontext connection")
	_ = client // use the client
}

func callNewClient()  {
	rawClient, err := rpc.Dial(infuraURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := ethclient.NewClient(rawClient)
	fmt.Println("We have a newclient connection")
	_ = client // use the client

}

func callBalanceAt()  {
	client, err := ethclient.Dial(infuraURL )
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	/*

	new(big.Float).SetInt(balance): 这里，big.Float是Go语言的一个类型，用于表示具有任意精度的浮点数。new(big.Float).SetInt(balance)这个操作会创建一个新的big.Float对象，并将balance（一个*big.Int类型，表示wei数量）设置为其值。

	new(big.Float).SetInt(big.NewInt(params.Ether)): 这是类似的操作，它创建了一个新的big.Float对象，并设置其值为params.Ether（这是一个常量，表示1 ether等于的wei数量）。

	new(big.Float).Quo(a, b): 这个操作执行了除法。Quo是big.Float类型的一个方法，它接受两个big.Float类型的参数，并返回它们的商。在这里，它将wei数量（balance）除以1 ether的wei数量（params.Ether），得到的结果就是ether数量。

	因此，整个表达式的意思是：将账户余额（以wei为单位）转换为ether。这是因为以太坊中的大多数操作都以wei为单位进行，但人们通常以ether为单位来理解和描述金额。

	`params` 是 Go Ethereum ("go-ethereum" 或 "geth") 包中的一个包。这个包包含了一些常量和函数，它们定义了 Ethereum 协议中的各种参数。

	比如，`params.Ether` 是一个常量，它定义了1 Ether等于多少Wei。Wei是以太坊网络中的最小货币单位，而Ether是最常用的单位。因此，这个常量常常被用于转换这两个单位。

	另外，`params` 包还定义了其他的一些参数，比如不同网络（Mainnet，Ropsten等）的链配置信息，Gas价格的一些默认值等等。这些都是和 Ethereum 协议以及其网络设置相关的参数。
	 */
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(big.NewInt(params.Ether)))
	fmt.Printf("Balance: %f ETH\n", ethValue) // output: Balance: XX.XXXXXX ETH

}


func callBlockByHash()  {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		panic(err)
	}

	transactionHash := common.HexToHash("0x32d18061d759dc52d5e2af4a96bc4b396034bac97d4e51a7f8fc600ff4cf7365") // replace with your transaction hash
	transaction, _, err := client.TransactionByHash(context.Background(), transactionHash)
	if err != nil {
		panic(err)
	}

	ethValue := new(big.Float).Quo(new(big.Float).SetInt(transaction.Value()), new(big.Float).SetInt(big.NewInt(params.Ether)))
	fmt.Printf("Transaction  : %f ETH \n", ethValue) // prints transaction value
}

func callBlockByNumber()  {
		client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/ovZAGPQlWJlYIoXzD1ifi4HDppTT4QS9")
		if err != nil {
			log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		}

		//如果你想获取最新的区块，你可以将 blockNumber 设置为 nil：
		//blockNumber := big.NewInt(17544277) // You can replace this with any block number you want to fetch.
		block, err := client.BlockByNumber(context.Background(), nil)
		if err != nil {
			log.Fatalf("Failed to get block: %v", err)
		}


	fmt.Println("Block number:", block.Number().Uint64())      // Block number
	fmt.Println("Block hash:", block.Hash().Hex())              // Block hash
	fmt.Println("Block timestamp:", convertBlockTimeToUTC8(block.Time()))
	fmt.Println("Block Nonce:", block.Nonce())                  // Block nonce
	fmt.Println("Number of transactions:", len(block.Transactions())) // Number of transactions in the block
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
	}



}

func callBlockNumber()  {
	 c := NewClient()
	blockNumber, err := c.client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get block number: %v", err)
	}

	fmt.Println("Current block number:", blockNumber)

}

func callCallContract ()  {

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	data, err := parsedABI.Pack("name")
	if err != nil {
		log.Fatalf("Failed to pack data for ABI: %v", err)
	}

	//列出 ABI中 所有的方法和参 数
	for name, method := range parsedABI.Methods {
		fmt.Printf("Method name: %s\n", name)
		fmt.Println("Input parameters:")
		for _, input := range method.Inputs {
			fmt.Printf("  name: %s, type: %s\n", input.Name, input.Type.String())
		}
		fmt.Println("Output parameters:")
		for _, output := range method.Outputs {
			fmt.Printf("  name: %s, type: %s\n", output.Name, output.Type.String())
		}
	}

	addr := common.HexToAddress(contractAddress)
	callMsg := ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call contract method: %v", err)
	}

	var name [32]byte
	err = parsedABI.UnpackIntoInterface(&name, "name", result)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}

	// Trim null characters
	trimmedName := strings.TrimRight(string(name[:]), "\x00")

	log.Printf("Contract name: %s", trimmedName)
}


func callCallContractAtHash() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}


	addr := common.HexToAddress(contractAddress)
	txHash := common.HexToHash("0x2d7877623e39fc55ab858f374d61870e4083d78a4e54b217c8c62a73f51d1236")
	blockHash, err := getBlockHashFromTxHash(txHash)
	if err != nil {
		log.Fatalf("Failed to get block hash from transaction hash: %v", err)
	}

	var (
		name string
		symbol string
		totalSupply *big.Int
		decimals uint8
		owner common.Address
	)

	// 创建一个查询每个常量的列表
	constants := []struct {
		name string
		val  interface{}
	}{
		{"name", &name},
		{"symbol", &symbol},
		{"totalSupply", &totalSupply},
		{"decimals", &decimals},
		{"owner", &owner},
	}



	for _, constant := range constants {
		// 打包方法调用以获得常量
		data, err := parsedABI.Pack(constant.name)
		if err != nil {
			log.Fatalf("Failed to pack data for %s: %v", constant.name, err)
		}

		callMsg := ethereum.CallMsg{
			To:   &addr,
			Data: data,
		}

		result, err := client.CallContractAtHash(context.Background(), callMsg, blockHash)
		if err != nil {
			log.Fatalf("Failed to call contract method: %v", err)
		}

		// 将结果解包到对应的变量中
		err = parsedABI.UnpackIntoInterface(constant.val, constant.name, result)
		if err != nil {
			log.Fatalf("Failed to unpack result for %s: %v", constant.name, err)
		}
	}

	// 打印结果
	// 打印结果
	log.Printf("name: %s", name)
	log.Printf("symbol: %s", symbol)
	log.Printf("totalSupply: %s", formatBigIntWithDecimals(totalSupply, decimals))
	log.Printf("decimals: %s", fmt.Sprintf("%d", decimals))
	log.Printf("owner: %s", owner.Hex())

}

func callChainID()  {
	c :=NewClient()

	chainID ,err := c.client.ChainID(context.Background())
	if err !=nil{
		panic(err)
	}

	fmt.Println("Chain ID:", chainID.String())

}


func callClient()  {
	c :=NewClient()

	rpcClient := c.client.Client()
	// 使用 rpcClient 进行底层的 RPC 操作
	// 例如发送原始的 JSON-RPC 请求
	_=rpcClient
	// 关闭客户端连接
	c.client.Close()
}

func callCodeAt()  {

	c :=NewClient()
	contractAddress := common.HexToAddress( contractAddress )
	code, err := c.client.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Contract bytecode: %x", code)

	c.client.Close()
}

func callEstimateGas()  {
	c :=NewClient()
	// 构造交易参数
	from := common.HexToAddress("0x2fcA988dBA4F3F2a3C5481450A5318b1484A077b")
	to := common.HexToAddress("0xCac0F1A06D3f02397Cfb6D7077321d73b504916e")
	amount := big.NewInt(1000000000000000000) // 1 ETH

	// 构造交易对象
	tx := ethereum.CallMsg{
		From:     from,
		To:       &to,
		GasPrice: nil, // 使用默认的燃料价格
		Value:    amount,
	}

	// 估算燃料消耗量
	gasLimit, err := c.client.EstimateGas(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Estimated gas limit: %d\n", gasLimit)
}

func callFeeHistory()  {

	c :=NewClient()
	// 获取最新的区块
	header, err := c.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// 最近的20个区块
	blockCount := uint64(20)

	// 奖励百分位数，比如50.0表示中位数
	rewardPercentiles := []float64{25.0, 50.0, 75.0}

	feeHistory, err := c.client.FeeHistory(context.Background(), blockCount, header.Number, rewardPercentiles)
	if err != nil {
		log.Fatal(err)
	}


	// 输出结果
	fmt.Printf("最老的区块: %s\n", feeHistory.OldestBlock.String())
	for i, baseFee := range feeHistory.BaseFee {
		// 创建新的big.Int实例，将i转换为big.Int类型后加到OldestBlock上
		blockNum := new(big.Int).Add(feeHistory.OldestBlock, big.NewInt(int64(i)))
		fmt.Printf("区块 %s: 基础费用: %s\n", blockNum.String(), baseFee.String())

		if i < len(feeHistory.Reward) { // 在这里添加检查，确保我们不超出范围
			for p, reward := range feeHistory.Reward[i] {
				if p < len(rewardPercentiles) {
					fmt.Printf("  %2.0f 百分位数矿工奖励: %s\n", rewardPercentiles[p], reward.String())
				}
			}
		}
		fmt.Println()
	}



}



type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}
const blockRange = 10

func callFilterLogs()  {
	c := NewClient()

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	address := common.HexToAddress(contractAddress)

	caller := bind.NewBoundContract(address, parsedABI, c.client, nil,nil )

	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}

	var result []interface{}

	if err := caller.Call(callOpts, &result, "name"); err != nil {
		log.Fatalf("Failed to retrieve token name: %v", err)
	}
	tokenName := result[0].(string)

	result = nil
	if err := caller.Call(callOpts, &result, "decimals"); err != nil {
		log.Fatalf("Failed to retrieve token decimals: %v", err)
	}
	decimals := result[0].(uint8)


	header, err := c.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve the latest block header: %v", err)
	}

	fromBlock := new(big.Int).Sub(header.Number, big.NewInt(blockRange))
	if fromBlock.Cmp(big.NewInt(0)) == -1 {
		fromBlock = big.NewInt(0)
	}

	transferEventSig := parsedABI.Events["Transfer"].ID

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   header.Number,
		Addresses: []common.Address{address},
	}
	ctx := context.Background()

	count := 0
	for count < 10 {
		logs, err := c.client.FilterLogs(ctx, query)
		if err != nil {
			log.Fatalf("Failed to filter logs: %v", err)
		}

		for _, vLog := range logs {
			if vLog.Topics[0].Hex() == transferEventSig.Hex() {
				var transferEvent TransferEvent
				err := parsedABI.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
				if err != nil {
					log.Fatalf("Failed to unpack data into Transfer event: %v", err)
				}

				transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
				transferEvent.Value = new(big.Int).SetBytes(vLog.Data)
				//ethValue := formatBigIntWithDecimals(transferEvent.Value,uint8(decimals))
				_=decimals
				blocktime := convertBlockTimeToUTC8(uint64(header.Time))
				txHash := vLog.TxHash.Hex() // 获取交易哈希

				isFromContract := checkIfContractAddress(c.client, transferEvent.From)
				isToContract := checkIfContractAddress(c.client, transferEvent.To)


				//判断地址是否是合约
				log.Printf("代币名称:%s|	交易时间: %s|	From: %s 是否为合约 : %t |	To: %s 是否为合约 : %t |	代币数量: %s |	交易哈希: %s \n",
					tokenName,
					blocktime,
					transferEvent.From.Hex(),
					isFromContract,
					transferEvent.To.Hex(),
					isToContract,
					transferEvent.Value,
					txHash)


				count++
				if count >= 10 {
					return
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func callHeaderByHash()  {
	c :=NewClient()
	hash 		:= common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	header, err := c.client.HeaderByHash(context.Background(), hash)
	if err != nil {
		log.Fatalf("Failed to get header by hash: %v", err)
	}

	fmt.Println("Block Number: ", header.Number.String()) // Block number
	fmt.Println("Time: ", header.Time) // Block time
	fmt.Println("TxHash: ", header.TxHash.Hex()) // Block
}

func callHeaderByNumber()  {
	c := NewClient()
	// 创建一个big.Int类型的变量表示区块号
	blockNumber := big.NewInt(17551002)

	header, err := c.client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("通过区块号获取区块头失败: %v", err)
	}

	fmt.Println(header.Number.String())     // 区块号
	fmt.Println(header.Time)                // 区块时间
	fmt.Println(header.Hash().Hex())        // 区块哈希
}

func callNetworkID()  {
	c :=NewClient()
	networkID, err := c.client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("获取网络ID失败: %v", err)
	}

	fmt.Println("Network ID:", networkID.String())
}

func callNonceAt()  {

	c :=NewClient()

	address :=common.HexToAddress("0x953cf65c1d08e29afa2daadb1133145bf0f64e99ae4358043e6de1f1e76c637c")
	nonce, err := c.client.NonceAt(context.Background(), address, nil) // for latest block
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	fmt.Printf("Nonce: %d\n", nonce)
}

func callPeerCount()  {
	c :=NewClient()
	peerCount, err := c.client.PeerCount(context.Background())
	if err != nil {
		log.Fatalf("Failed to get peer count: %v", err)
	}

	fmt.Printf("Peer count: %d\n", peerCount)
}

func callPendingBalanceAt()  {
	c :=NewClient()
	address := common.HexToAddress("0x2d0c46c14934af9773755e53bb2c7237b0b4e80e")
	pendingBalance, err := c.client.PendingBalanceAt(context.Background(), address)
	if err != nil {
		log.Fatalf("Failed to get pending balance: %v", err)
	}

	ethvalue := formatBigIntWithDecimals(pendingBalance,18)
	fmt.Printf("Pending balance: %s\n", ethvalue)

}

func callPendingCallContract()  {
	 c:=NewClient()
	ctx := context.Background()
	address := common.HexToAddress(contractAddress )

	code, err := c.client.PendingCodeAt(ctx, address)
	if err != nil {
		log.Fatalf("Failed to get pending code at address: %v", err)
	}

	fmt.Printf("Pending code: %s\n", code)
}

func callPendingCodeAt()  {
	c:=NewClient()
	ctx := context.Background()
	address := common.HexToAddress(contractAddress )

	code, err := c.client.PendingCodeAt(ctx, address)
	if err != nil {
		log.Fatalf("Failed to get pending code at address: %v", err)
	}

	fmt.Printf("Pending codeat: %d\n", code)
}

func callPendingNonceAt()  {
	c :=NewClient()
	account := common.HexToAddress(oeaAddress)
	nonce, err := c.client.PendingNonceAt(context.Background(), account)
	if err != nil {
		log.Fatalf("Failed to get pending nonce: %v", err)
	}

	log.Printf("Pending nonce: %v", nonce)
}

func callPendingStorageAt()  {
	 c :=NewClient()
	contractAddress := common.HexToAddress(contractAddress)
	key := common.BigToHash(big.NewInt(1)) // 用于获取第一个状态变量的值


	value, err := c.client.PendingStorageAt(context.Background(), contractAddress, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)  //
}

func callPendingTransactionCount()  {
	c :=NewClient()
	nonce, err := c.client.PendingTransactionCount(context.Background())
	if err != nil {
		log.Fatal(err)
	}
//在这个例子中，PendingTransactionCount方法返回当前以太坊节点的事务池中正在等待的事务数量。这些事务正在等待被矿工打包进一个新的区块中。
	fmt.Printf("pending tx count: %d\n", nonce)
}

func callSendTransaction()  {

	//eth 测试链
	client, err := ethclient.Dial(ethGoerliRpc)
	if err != nil {
		log.Fatal(err)
	}


	privateKey, err := crypto.HexToECDSA("priavte_key")
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("获取Nonce失败: %v", err)
	}

	toAddress := common.HexToAddress("address")

	value := big.NewInt(10000000000000000) // 0.01 eth
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}
	var data []byte
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get network ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign tx: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent: %s", signedTx.Hash().Hex()) // print transaction hash
}

func callStorageAt()  {
	c :=NewClient()
	contractAddress := common.HexToAddress(contractAddress )
	//blockNumber := big.NewInt(17553839) // change this with your actual block number

	// 创建等待组和结果通道
	var wg sync.WaitGroup
	results := make(chan string)

	// 设置并发的 goroutine 数量
	concurrency := 10

	// 分配任务给多个 goroutine
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(startSlot, endSlot int) {
			defer wg.Done()
			for slot := startSlot; slot <= endSlot; slot++ {
				storageKey := common.BigToHash(big.NewInt(int64(slot)))
				result, err := c.client.StorageAt(context.Background(), contractAddress, storageKey, nil)
				if err != nil {
					log.Println(err)
					continue
				}
				value := common.Bytes2Hex(result)
				results <- fmt.Sprintf("Storage slot %d: %s", slot, value)
			}
		}(i*1000, (i+1)*1000-1)
	}

	// 等待所有任务完成并关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 读取结果通道并打印结果
	for result := range results {
		fmt.Println(result)
	}

}

func callSubscribeFilterLogs()  {

	c :=NewClient()
	// 合约地址和 ABI
	contractAddress := common.HexToAddress(contractAddress)
	// 从 ABI 中解析出 Transfer 事件的签名
	contractABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	eventSignature := contractABI.Events["Transfer"].ID

	// 创建日志订阅器
	logs := make(chan types.Log)
	sub, err := c.client.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignature}},
	}, logs)
	if err != nil {
		log.Fatal(err)
	}

	// 处理接收到的日志
	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case log := <-logs:
				from := common.HexToAddress(log.Topics[1].Hex())
				to := common.HexToAddress(log.Topics[2].Hex())
				value := new(big.Int)
				value.SetBytes(log.Data)
				fmt.Printf("Transfer event: from=%s, to=%s, value=%s\n", from.Hex(), to.Hex(), value.String())
			}
		}
	}()

	// 等待程序退出
	select {}
}

func callSubscribeNewHead()  {
	 c:=NewClient()
	headers := make(chan *types.Header )
	sub, err := c.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Printf("New block header: %v\n", header)
			// 在这里处理新区块头的逻辑
		}
	}
}

func callSuggestGasPrice()  {
	c :=NewClient()
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gwei := new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9))
	fmt.Printf("建议燃气价格: %s Gwei\n", gwei.Text('f', 2))
}

func callSuggestGasTipCap()  {
	 c:=NewClient()
	// 获取建议的燃气小费上限
	tipCap, err := c.client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 将 Wei 转换为 Gwei
	gweiValue := new(big.Float).Quo(new(big.Float).SetInt(tipCap), big.NewFloat(1e9))

	fmt.Println("建议的燃气小费上限:", gweiValue.String(), "Gwei")
}

func callcallSyncProgress()  {
	 c :=NewClient()
	// 获取同步进度
	// 获取同步状态
	prog, err := c.client.SyncProgress(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 打印同步进度信息
	// 打印同步进度信息
	log.Printf("当前区块数: %v", prog.CurrentBlock)
	log.Printf("最高区块数: %v", prog.HighestBlock)

	// 计算同步百分比
	progressPercent := new(big.Float).Quo(
		new(big.Float).Mul(
			new(big.Float).SetUint64(prog.CurrentBlock),
			big.NewFloat(100),
		),
		new(big.Float).SetUint64(prog.HighestBlock),
	)
	log.Printf("同步进度: %.2f%%", progressPercent)
}

func callTransactionByHash()  {
	 c:=NewClient()
	// 交易哈希
	txHash := common.HexToHash("0x571f660a8f0269cdb188d28ea886486315f10b34403e8783c79fce3a9349363c")

	// 调用方法获取交易
	tx, isPending, err := c.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	// 打印交易信息
	if isPending {
		fmt.Println("交易还在待定状态")
	} else {
		fmt.Println("交易已经被确认")
	}



	// 打印交易信息
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("是否待定: %v\n", isPending)
	fmt.Printf("交易接收者: %s\n", tx.To().Hex())
	fmt.Printf("交易金额: %s\n", tx.Value().String())
	fmt.Printf("交易数据: %s\n", string(tx.Data()))
}

func callTransactionCount()  {
	 c :=NewClient()
	header, err := c.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	blockHash := header.Hash()

	count, err := c.client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction count: %d\n", count)
}

func callTransactionInBlock()  {
	 c :=NewClient()

	header, err := c.client.HeaderByNumber(context.Background(),nil)

	// 要查询的块哈希和事务索引
	index := uint(2)

	// 获取指定块中的事务
	tx, err := c.client.TransactionInBlock(context.Background(), header.Hash(), index)
	if err != nil {
		log.Fatal(err)
	}

	from, err := c.client.TransactionSender(context.Background(), tx, header.Hash(), index)
	if err != nil {
		log.Fatal(err)
	}

	// 打印事务信息
	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("Transaction From: %s\n", from.Hex())
	fmt.Printf("Transaction To: %s\n", tx.To().Hex())
	fmt.Printf("Transaction Value: %s\n", tx.Value().String())
	fmt.Printf("Transaction Data: %s\n", string(tx.Data()))
}

func callTransactionReceipt()  {
	c :=NewClient()
	// 交易哈希
	txHash := common.HexToHash("0xfd3f1c07cd7587006bbaa5c82149df99290778c1aba681256114e5d4ef1b1c85")

	// 获取交易收据
	receipt, err := c.client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	// 输出交易相关信息
	fmt.Println("交易哈希:", receipt.TxHash.Hex())
	fmt.Println("交易索引:", receipt.TransactionIndex)
	fmt.Println("区块哈希:", receipt.BlockHash.Hex())
	fmt.Println("区块号:", receipt.BlockNumber)
	fmt.Println("燃气消耗:", receipt.GasUsed)
	fmt.Println("合约地址:", receipt.ContractAddress.Hex())
	fmt.Println("日志数量:", len(receipt.Logs))

}

func callTransactionSender()  {

	c :=NewClient()
	txHash := common.HexToHash("0x571f660a8f0269cdb188d28ea886486315f10b34403e8783c79fce3a9349363c")
	tx, isPending, err := c.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	if isPending {
		fmt.Println("Transaction is pending")
	} else {
		receipt, err := c.client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			log.Fatal(err)
		}

		block, err := c.client.BlockByHash(context.Background(), receipt.BlockHash)
		if err != nil {
			log.Fatal(err)
		}

		sender, err := c.client.TransactionSender(context.Background(), tx, receipt.BlockHash, uint(receipt.TransactionIndex))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("交易已包含在区块中:")
		fmt.Println("区块哈希:", receipt.BlockHash.Hex())
		fmt.Println("区块编号:", block.Number().Uint64())
		blockTime := time.Unix(int64(block.Time()), 0)
		fmt.Println("区块时间戳:", blockTime.Format(time.RFC3339))
		fmt.Println("交易发送者:", sender.Hex())
	}
}