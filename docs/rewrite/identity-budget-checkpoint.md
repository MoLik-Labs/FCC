# Current Chat identity budget checkpoint

The current Python Chat identity resolver can share a source-output retention
budget. The Go adapter now accepts an optional pending budget for
`ResolveCurrent`; the historical `Resolve` behavior remains unchanged. The Go
budget models byte and record limits separately, leaves both counters unchanged
on a rejected addition, and provides release accounting. Identity allocation
charges 32 bytes and one record, while a newly learned ID charges its UTF-8
byte length without adding a record. As in Python, if that latter charge fails
a newly allocated slot and its index remain present.

`bench_local/generate_golden_identity_budget.py` executes the current Python
resolver and writes ten independent reference scenarios. These cover repeated
IDs, exact boundaries, failed allocation, failed ID charges, multibyte IDs,
record limits, ambiguity and subsequent steps. Go replays each step and
compares returned values, error text, budget counters and identity state. The
fixture SHA-256 is
`608eaae70d48489c977f6d27262146e5e2ef9b65ad3aa19200a6fb5e49ef0fb8`.
Go vet, tests and race tests passed; 71 related Python tests passed.

A local Go microbenchmark creating and resolving eight index-only identities
reported 671–745 ns and 432 B per operation without a budget versus 732–780 ns
and 464 B with a budget, over three short runs. This includes constructing the
budget and excludes ID-byte charging; it is not a full-stream performance claim.

This is only an identity-accounting slice. The Go counter uses machine integers
and assumes representable values; Python has arbitrary-precision integers.
UTF-8 ID charging assumes valid scalar text, and the Go error currently mirrors
the message rather than all Python failure metadata. The standalone
source-order buffer is now covered separately in the ordered-buffer checkpoint;
the publication path and current provider adapter remain to be ported and
checked against Python before release.
