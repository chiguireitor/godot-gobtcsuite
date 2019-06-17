extends Spatial

const GoBtcSuite = preload("res://bin/gobtcsuite.gdns")
const wallet_filename = "user://wallet.dat"
const SERVERS = {
	testnet="http://localhost:14100"
	}
const UNSPENDABLE = {
	testnet="mvCounterpartyXXXXXXXXXXXXXXW24Hef",
	mainnet="1CounterpartyXXXXXXXXXXXXXXXUWLpVr"
	}

var burn_address
var counterblock_server

func open_wallet(network):
	burn_address = UNSPENDABLE[network]
	counterblock_server = SERVERS[network]
	var saved_wallet = File.new()
	if not saved_wallet.file_exists(wallet_filename):
		var Btc = GoBtcSuite.new()
		Btc.set_network(network)
		var masterseed = Btc.gen_masterseed()
		
		saved_wallet.open(wallet_filename, File.WRITE)
		saved_wallet.store_line(masterseed)
		
		return Btc
	else:
		saved_wallet.open(wallet_filename, File.READ)
		var masterseed = saved_wallet.get_line()
		var Btc = GoBtcSuite.new()
		Btc.set_network(network)
		Btc.load_masterseed(masterseed)
		
		return Btc
		
	
func _ready():
	var Btc = open_wallet("testnet")
	var addr = Btc.get_address()
	#var msg = "This is a test"
	#var signed = Btc.sign_message(msg)
	print(addr)
	#print(signed)
	#print(Btc.verify_message(addr, msg, signed))
	#http.connect("request_completed", self, "_on_HTTPRequest_request_completed")
	counterblock("get_normalized_balances", {addresses=[addr]}, funcref(self, "process_balances"))
	#counterparty("get_block_info", {block_index=581010}, funcref(self, "process_balances"))
	counterparty("create_issuance", {source=addr, quantity=1, divisible=false, description="TESTSWORD", asset="A26936834624"}, funcref(self, "process_balances"))
	
func process_balances(res):
	print(res)
	
var _id_seq_ = 0
var seq_cbs = {}
func counterparty(method, params, fref):
	counterblock("proxy_to_counterpartyd", {method=method, params=params}, fref)
	
func counterblock(method, params, fref):
	var rdata = JSON.print({
		jsonrpc="2.0",
		id=_id_seq_,
		method=method,
		params=params
	})
	var headers = [
		"Content-type: application/json"
	]
	var req = HTTPRequest.new()
	req.use_threads = true
	add_child(req)
	req.connect("request_completed", self, "_on_HTTPRequest_request_completed")
	req.request(counterblock_server, headers, true, HTTPClient.METHOD_POST, rdata)
	seq_cbs[_id_seq_] = {fref=fref, req=req}
	_id_seq_ += 1
	
func _on_HTTPRequest_request_completed( result, response_code, headers, body ):
	var json = JSON.parse(body.get_string_from_utf8())
	var id = int(json.result["id"])
	
	if seq_cbs.has(id):
		var cb = seq_cbs[id]
		seq_cbs.erase(id)
		cb["fref"].call_func(json.result)
		cb["req"].queue_free()


