#!/usr/bin/env bash
chmod +x ./ab-srv
chmod +x ./ab-client

# ./ab-srv>/dev/null&
./ab-srv
SRV_PID=$!

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
    ./ab-client l User1 Pass1 192.168.1.1 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client l User1 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client r User1 Pass1 192.168.1.1 >/tmp/ab.out

./ab-client l User1 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

./ab-client ab 192.168.1.0/25 >/tmp/ab.out

./ab-client l User2 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client db 192.168.1.0/25 >/tmp/ab.out

./ab-client l User2 Pass1 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

for ((i=0; i < 100; i++))
do
    ./ab-client l User_$i Pass2 192.168.2.1 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client l User$i Pass2 192.168.2.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client aw 192.168.2.0/25 >/tmp/ab.out

./ab-client l User_100 Pass2 192.168.2.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

./ab-client dw 192.168.1.1 >/tmp/ab.out

./ab-client l User_100 Pass2 192.168.1.1 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client r User Pass2 192.168.1.2 >/tmp/ab.out

./ab-client l User Pass2 192.168.1.2 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

for ((i=0; i < 1000; i++))
do
    ./ab-client l User_$i Pass_$i 192.168.1.3 >/tmp/ab.out
    fileEquals /tmp/ab.out "${expected_true}"
done
./ab-client l User_1000 Pass_1000 192.168.1.3 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_false}"

./ab-client r User_1000 Pass_1000 192.168.1.3 >/tmp/ab.out

./ab-client l User_1000 Pass_1000 192.168.1.3 >/tmp/ab.out
fileEquals /tmp/ab.out "${expected_true}"

kill ${SRV_PID} 2>/dev/null || true

rm -f /tmp/ab.out
echo "PASS"