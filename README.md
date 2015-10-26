# T-Digest

This is a reimplementation of [t-digest][1] in Go that uses amortized merging
instead of trees. It provides `MergingDigest` type that is similar to the
original `MergingDigest`, however it does not have bounds on the summary size,
and is more akin to tree digests in that regard.

WARNING: this is alpha and API is subject to change

[1]: https://github.com/tdunning/t-digest
