
runFile="outMultiple.lexe"
# kill all screens
killall screen
# restart the program
screen -S targetSession -d -m
screen -r targetSession -X stuff "./$runFile"$(echo -ne '\015')

sleep 1

# test using pgrep
if ! pgrep "out" > /dev/null
then
    echo "no program started"
else
    echo "program starts again"
fi