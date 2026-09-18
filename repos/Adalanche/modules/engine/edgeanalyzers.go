package engine

import (
	"sync"

	"github.com/lkarlslund/adalanche/modules/ui"
)

//go:generate go tool github.com/dmarkham/enumer -type=ProcessPriority -output enums.go

type ProgressCallbackFunc func(progress int, totalprogress int)

type ProcessorFunc func(ao *IndexedGraph)
type ReadOnlyProcessorFunc func(view *FrozenGraph)
type NodePatchProcessorFunc func(view *FrozenGraph, out *NodePatchSet)
type EdgeDeltaProcessorFunc func(view *FrozenGraph, out *EdgeDelta)

type ProcessPriority int

const (
	BeforeMergeLow ProcessPriority = iota
	BeforeMerge
	BeforeMergeHigh
	BeforeMergeFinal
	AfterMergeLow
	AfterMerge
	AfterMergeHigh
	AfterMergeFinal
)

type ProcessorKind int

const (
	ProcessorKindGraphMutator ProcessorKind = iota
	ProcessorKindReadOnly
	ProcessorKindNodePatch
	ProcessorKindEdgeDelta
)

type processorInfo struct {
	description string
	priority    ProcessPriority
	loader      LoaderID
	kind        ProcessorKind
	mutator     ProcessorFunc
	readOnly    ReadOnlyProcessorFunc
	nodePatch   NodePatchProcessorFunc
	edgeDelta   EdgeDeltaProcessorFunc
}

var registeredProcessors []processorInfo

func (l LoaderID) AddProcessor(pf ProcessorFunc, description string, priority ProcessPriority) {
	l.AddGraphMutator(pf, description, priority)
}

func (l LoaderID) AddGraphMutator(pf ProcessorFunc, description string, priority ProcessPriority) {
	registeredProcessors = append(registeredProcessors, processorInfo{
		loader:      l,
		description: description,
		priority:    priority,
		kind:        ProcessorKindGraphMutator,
		mutator:     pf,
	})
}

func (l LoaderID) AddReadOnlyProcessor(pf ReadOnlyProcessorFunc, description string, priority ProcessPriority) {
	registeredProcessors = append(registeredProcessors, processorInfo{
		loader:      l,
		description: description,
		priority:    priority,
		kind:        ProcessorKindReadOnly,
		readOnly:    pf,
	})
}

func (l LoaderID) AddNodePatchProcessor(pf NodePatchProcessorFunc, description string, priority ProcessPriority) {
	registeredProcessors = append(registeredProcessors, processorInfo{
		loader:      l,
		description: description,
		priority:    priority,
		kind:        ProcessorKindNodePatch,
		nodePatch:   pf,
	})
}

func (l LoaderID) AddEdgeDeltaProcessor(pf EdgeDeltaProcessorFunc, description string, priority ProcessPriority) {
	registeredProcessors = append(registeredProcessors, processorInfo{
		loader:      l,
		description: description,
		priority:    priority,
		kind:        ProcessorKindEdgeDelta,
		edgeDelta:   pf,
	})
}

// LoaderID = wildcard
func Process(ao *IndexedGraph, statustext string, l LoaderID, priority ProcessPriority) error {
	var priorityProcessors []processorInfo
	for _, potentialProcessor := range registeredProcessors {
		if (potentialProcessor.loader == l || l == -1) && potentialProcessor.priority == priority {
			priorityProcessors = append(priorityProcessors, potentialProcessor)
		}
	}

	aoLen := ao.Order()
	total := len(priorityProcessors) * aoLen

	if total == 0 {
		return nil
	}

	// We need to process this many objects
	pb := ui.ProgressBar(statustext, int64(total))
	var analysisProcessors []processorInfo
	for _, processor := range priorityProcessors {
		if processor.kind == ProcessorKindGraphMutator {
			processor.mutator(ao)
			pb.Add(int64(aoLen))
			continue
		}
		analysisProcessors = append(analysisProcessors, processor)
	}

	if len(analysisProcessors) > 0 {
		view := ao.Freeze()
		nodePatchResults := make([]NodePatchSet, len(analysisProcessors))
		edgeDeltaResults := make([]EdgeDelta, len(analysisProcessors))

		var wg sync.WaitGroup
		for i, processor := range analysisProcessors {
			wg.Add(1)
			go func(i int, processor processorInfo) {
				defer wg.Done()

				switch processor.kind {
				case ProcessorKindReadOnly:
					processor.readOnly(view)
				case ProcessorKindNodePatch:
					processor.nodePatch(view, &nodePatchResults[i])
				case ProcessorKindEdgeDelta:
					processor.edgeDelta(view, &edgeDeltaResults[i])
				}

				pb.Add(int64(aoLen))
			}(i, processor)
		}
		wg.Wait()

		var dropIndexes bool
		for i := range nodePatchResults {
			nodePatchResults[i].Apply(ao)
			dropIndexes = dropIndexes || nodePatchResults[i].HasOperations()
		}
		if dropIndexes {
			ao.DropIndexes()
		}

		for i := range edgeDeltaResults {
			edgeDeltaResults[i].Apply(ao)
		}
	}

	pb.Finish()

	return nil
}
