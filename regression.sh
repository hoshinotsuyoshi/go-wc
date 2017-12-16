#!/bin/sh -e
if [ -n "$CI" ]; then
  set -x
fi

ORIG_WC=/usr/bin/wc
MY_WC=$PWD/go-wc

cd /tmp

## regression 1. stdin

echo "xxx" | $ORIG_WC 2>orig_stderr 1>orig_stdout
echo $? >orig_exitcode
echo "xxx" | $MY_WC 2>my_stderr 1>my_stdout
echo $? >my_exitcode

diff -u orig_stderr   my_stderr
# diff -u orig_stdout   my_stdout
diff -u orig_exitcode my_exitcode


## regression 2. default option

echo "あああ\n" > /tmp/x
echo "bbb ccc\nいいい\n" > /tmp/y

$ORIG_WC -- /tmp/x ./y 2>orig_stderr 1>orig_stdout
echo $? >orig_exitcode
$ORIG_WC -- /tmp/x ./y 2>my_stderr 1>my_stdout
echo $? >my_exitcode

diff -u orig_stderr   my_stderr
diff -u orig_stdout   my_stdout
diff -u orig_exitcode my_exitcode


## regression 3. m option

echo "あああ\n" > /tmp/x
echo "bbb ccc\nいいい\n" > /tmp/y

$ORIG_WC -m -- /tmp/x ./y 2>orig_stderr 1>orig_stdout
echo $? >orig_exitcode
$ORIG_WC -m -- /tmp/x ./y 2>my_stderr 1>my_stdout
echo $? >my_exitcode

diff -u orig_stderr   my_stderr
diff -u orig_stdout   my_stdout
diff -u orig_exitcode my_exitcode
