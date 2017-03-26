require 'socket'

sock = TCPSocket.new("localhost", 8080)
sock.write("\x01\x01\x00")



sock = TCPSocket.new("localhost", 8080)
sock.write("\x05\x00")

sock = TCPSocket.new("localhost", 8080)
sock.write("\x05\x01\x00")

sock = TCPSocket.new("localhost", 8080)
sock.write("\x05\x02\x00\x01")

sock = TCPSocket.new("localhost", 8080)
sock.write("\x05\x01\x01")
