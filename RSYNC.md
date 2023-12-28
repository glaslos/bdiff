## The rsync algorithm 

[Source](https://rsync.samba.org/tech_report/node2.html)

Suppose we have two general purpose computers $\alpha$ and $\beta$. Computer $\alpha$ has access to a file A and $\beta$ has access to file B, where A and B are ``similar''. There is a slow communications link between $\alpha$ and $\beta$.

The rsync algorithm consists of the following steps:

1. $\beta$ splits the file B into a series of non-overlapping fixed-sized blocks of size S bytes1. The last block may be shorter than S bytes.

2. For each of these blocks $\beta$ calculates two checksums: a weak ``rolling'' 32-bit checksum (described below) and a strong 128-bit MD4 checksum.

3. $\beta$ sends these checksums to $\alpha$.

4. $\alpha$ searches through A to find all blocks of length S bytes (at any offset, not just multiples of S) that have the same weak and strong checksum as one of the blocks of B. This can be done in a single pass very quickly using a special property of the rolling checksum described below.

5. $\alpha$ sends $\beta$ a sequence of instructions for constructing a copy of A. Each instruction is either a reference to a block of B, or literal data. Literal data is sent only for those sections of A which did not match any of the blocks of B.

The end result is that $\beta$ gets a copy of A, but only the pieces of A that are not found in B (plus a small amount of data for checksums and block indexes) are sent over the link. The algorithm also only requires one round trip, which minimises the impact of the link latency.

The most important details of the algorithm are the rolling checksum and the associated multi-alternate search mechanism which allows the all-offsets checksum search to proceed very quickly. These will be discussed in greater detail below.
