#!/bin/env perl

use v5.38;
use utf8;
use strict;
use FindBin;

# f16_lt.txt is generated by TestFloat-3b/testfloat_gen.
# http://www.jhauser.us/arithmetic/TestFloat.html
# $ ./testfloat_gen -level 1 f16_lt > f16_lt.txt
open my $fh, "<", "$FindBin::Bin/f16_lt.txt" or die "Can't open f16_lt.txt: $!";

say <<EOF;
// Code generated by scripts/f16_lt.pl; DO NOT EDIT.

package float16

var f16Lt = []struct {
    a, b Float16
    want bool
} {
EOF
while(my $line = <$fh>) {
    chomp $line;
    my ($a, $b, $c, undef) = split /\s+/, $line;
    my $want = $c ne "0" ? "true" : "false";
    say "    {0x$a, 0x$b, $want},";
}

say "}";

close $fh;
