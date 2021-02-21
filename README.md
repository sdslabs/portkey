<h1 align="center">Portkey</h1>
 <p align="center"><b>A peer to peer file transfer tool using RTC over QUIC</b></p>

#
Portkey is a file transfer tool that uses RTC over QUIC protocol to provide peer to peer file transfers. It uses the zstd on-the-fly compression algorithm alongside with QUIC protocol to achieve high speed file transfers. 

## Build 

To build portkey from source, first git clone the repo 

        git clone git@github.com:sdslabs/portkey.git

Portkey strictly uses go 1.14. All dependencies are vendored so no need to install any dependencies. Use the makefile to build the portkey binary. Simply run

        make build 

This will place the portkey binary in your working directory

## Usage
 
You can use the following flags

        -b, --benchmark      Set to benchmark locally(for local testing)
        -h, --help           Help for portkey
        -k, --key string     Key to connect to peer
        -r, --receive        Set to receive files
        -p, --rpath string   Absolute path of where to receive files, pwd by default
        -s, --send string    Absolute path of directory/file to send

If you wish to send a file to someone, initiate the transfer by using the -s flag and specifying the path to the file you want to send

        portkey -s /example/path/to/file

You will then receive a key. Pass that key to the recipient. They can now use 

        portkey -r -k example_key -p /example/destination/path

If -p is not specified, the file will be received in the present working directory. You can also specify both -r and -s if you wish to send and receive files simultaneously. Don't use -r if you don't intend on receiving files.

## Note
If after the key exchange, portkey is stuck at 

        starting ice connection... 

then the file transfer is unfortunately not possible because you or your peer are either behind a firewall or a misbehaved NAT that makes a peer to peer connection impossible. You can try turning off your firewall or shifting to an un-NAT-ed connection like your mobile data if you are on a wi-fi connection. 
## License

Licensed under the [MIT License](./LICENSE).