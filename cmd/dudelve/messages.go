package main

type ScanProgress struct {
	scannedCount uint64
	scannedOk    bool
	moreCounts   chan uint64
	err          chan error
	result       chan *node
}

func (sp ScanProgress) awaitNext() ScanProgress {
	scanned, ok := <-sp.moreCounts
	return ScanProgress{scannedCount: scanned, scannedOk: ok, moreCounts: sp.moreCounts, err: sp.err, result: sp.result}
}
