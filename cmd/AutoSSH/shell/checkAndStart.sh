runFile2="outMultiple.lexe"

if ! pgrep "out" > /dev/null
then
    if test -f $runFile2
    then
        echo "start the program"
        killall screen
		screen -S targetSession -d -m
		screen -r targetSession -X stuff "./$runFile2"$(echo -ne '\015')
    fi
else
    echo "program is running"
fi