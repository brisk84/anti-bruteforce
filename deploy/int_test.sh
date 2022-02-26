#!/usr/bin/env bash
chmod +x ./ab-srv
chmod +x ./ab-client

./ab-srv -testMode=true>/dev/null&
SRV_PID=$!
sleep 1

expected_true='ok=true'
expected_false='ok=false'

fileEquals()
{
    local fileData
    fileData=$(cat "$1")
    if [[ "$fileData" != ${2} ]]; then
        echo -e "unexpected output, $1:\n${fileData}"
        kill ${SRV_PID} 2>/dev/null || true
        exit 1
    fi
}

for ((i=0; i < 10; i++))
do
    ./ab-client login User1 Pass1 192.168.1.1 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client login User1 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

sleep 10

./ab-client login User1 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

for ((i=0; i < 10; i++))
do
    ./ab-client login User2 Pass1 192.168.1.1 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client login User2 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client reset User2 Pass1 192.168.1.1 >/tmp/ab.out

./ab-client login User2 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

./ab-client addBlack 192.168.1.0/25 >/tmp/ab.out

./ab-client login User2 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client delBlack 192.168.1.0/25 >/tmp/ab.out

./ab-client login User2 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

for ((i=0; i < 100; i++))
do
    ./ab-client login User_$i Pass2 192.168.2.1 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client login User$i Pass2 192.168.2.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client addWhite 192.168.2.0/25 >/tmp/ab.out

./ab-client login User_100 Pass2 192.168.2.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

./ab-client delWhite 192.168.1.1 >/tmp/ab.out

./ab-client login User_100 Pass2 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client reset User Pass2 192.168.1.2 >/tmp/ab.out

./ab-client login User Pass2 192.168.1.2 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

for ((i=0; i < 1000; i++))
do
    ./ab-client login User_$i Pass_$i 192.168.1.3 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client login User_1000 Pass_1000 192.168.1.3 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

sleep 10

./ab-client login User_1000 Pass_1000 192.168.1.3 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

kill ${SRV_PID} 2>/dev/null || true

rm -f /tmp/ab.out
echo "PASS"