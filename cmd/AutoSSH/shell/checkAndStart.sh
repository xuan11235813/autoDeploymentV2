runFile1="tcpKeepAliveServer"
runDirectory1="radarSystem/remotePort"
runFile2="netMultiple.lexe"
runDirectory2="radarSystem/radarSystem"

if ! pgrep "netMultiple" > /dev/null
then
    if test -f $runDirectory2/$runFile2
    then
        echo "start the program"
        killall screen
		screen -S targetSessionR -d -m
        screen -r targetSessionR -X stuff "cd $runDirectory2"$(echo -ne '\015')
		screen -r targetSessionR -X stuff "./$runFile2"$(echo -ne '\015')
        screen -S targetSessionP -d -m
        screen -r targetSessionP -X stuff "cd $runDirectory1"$(echo -ne '\015')
		screen -r targetSessionP -X stuff "./$runFile1"$(echo -ne '\015')
    fi
else
    echo "program is running"
fi