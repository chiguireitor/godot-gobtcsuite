extends Spatial

const GoBtcSuite = preload("res://bin/gobtcsuite.gdns")
const wallet_filename = "user://wallet.dat"

func open_wallet(network):
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
	var Btc = open_wallet("mainnet")
	print(Btc.get_address())
	print(Btc.sign_message("This is a test"))


