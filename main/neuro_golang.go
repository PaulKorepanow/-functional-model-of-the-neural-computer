package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Neuron base element arhitecture
type Neuron struct {
	weight []float64
	input  []float64
	output float64
}

func (n *Neuron) transferFunction(input []float64) {
	n.input = input
	n.output = 0
	for i := 0; i <= len(n.input)-1; i++ {
		n.output += n.input[i] * n.weight[i]
	}
	n.output = n.sigmoid()
}

func (n Neuron) sigmoid() float64 {
	return 1 / (1 + math.Exp(-n.output))
}

//NeuralNetwork with all attributes
type NeuralNetwork struct {
	learningRate    float64
	numNeurons      []int
	archNeuronNet   [][]Neuron
	hiddenOut       []float64
	outputOut       []float64
	weightDelta     [][]float64
	err             []float64
	errback         float64
	correctAnswer   int
	incorrectAnswer int
}

func (p *NeuralNetwork) init() {
	p.archNeuronNet = make([][]Neuron, len(p.numNeurons)-1)
	neurons := p.numNeurons[1:]
	weights := p.numNeurons[:2]
	for i := range p.archNeuronNet {
		p.archNeuronNet[i] = make([]Neuron, 0)
	}
	for j, ja := range neurons {
		p.archNeuronNet[j] = make([]Neuron, ja)
	}
	for ia := range p.archNeuronNet {
		for j, ja := range p.archNeuronNet[ia] {
			for ind := 0; ind <= weights[ia]-1; ind++ {
				val := rand.Float64()
				ja.weight = append(ja.weight, val)
			}
			p.archNeuronNet[ia][j] = ja
		}
	}

	for ind := 0; ind <= neurons[1]; ind++ {
		p.err = make([]float64, neurons[1])
		val := rand.Float64()
		p.err = append(p.err, val)
	}

}

type dataStream interface {
	dataStreamRight(inputVals []int, expectedVals []int, b bool)
}

func (p *NeuralNetwork) dataStreamRight(inputVals []float64, expectedVals []float64, edu bool) {
	for i := 0; i <= p.numNeurons[1]-1; i++ {
		p.archNeuronNet[0][i].transferFunction(inputVals)
		p.hiddenOut = append(p.hiddenOut, p.archNeuronNet[0][i].output)
	}
	for i := 0; i <= p.numNeurons[2]-1; i++ {
		p.archNeuronNet[1][i].transferFunction(p.hiddenOut)
		p.outputOut = append(p.outputOut, p.archNeuronNet[1][i].output)
	}
	p.hiddenOut = nil
	if edu {
		p.outputOut = nil
		p.dataStreamBack(inputVals, expectedVals)
	} else {
		var compereVal1 float64
		var compereVal2 float64
		index1 := 0
		index2 := 0
		for i1, v1 := range expectedVals {
			if v1 > compereVal1 {
				compereVal1 = v1
				index2 = i1

			}
		}
		for i2, v2 := range p.outputOut {
			if v2 > compereVal2 {
				compereVal2 = v2
				index1 = i2
			}
		}
		if index1 == index2 {
			p.correctAnswer++
		} else {
			p.incorrectAnswer++
		}
	}
	p.outputOut = p.outputOut[:0]

}

func (p *NeuralNetwork) dataStreamBack(input []float64, expectedVals []float64) {
	p.weightDelta = make([][]float64, 2)
	for neuronOut := 0; neuronOut <= p.numNeurons[2]-1; neuronOut++ {
		p.err[neuronOut] = p.archNeuronNet[1][neuronOut].output - expectedVals[neuronOut]
		delta := p.err[neuronOut] * p.archNeuronNet[1][neuronOut].output * (1 - p.archNeuronNet[1][neuronOut].output)
		p.weightDelta[0] = append(p.weightDelta[0], delta)

		for neuronHid := 0; neuronHid <= p.numNeurons[1]-1; neuronHid++ {
			p.archNeuronNet[1][neuronOut].weight[neuronHid] = p.archNeuronNet[1][neuronOut].weight[neuronHid] -
				p.archNeuronNet[0][neuronHid].output*
					p.learningRate*
					delta

		}
	}
	for neuronHid := 0; neuronHid <= p.numNeurons[1]-1; neuronHid++ {
		p.errback = 0
		for neuronOut := 0; neuronOut <= p.numNeurons[2]-1; neuronOut++ {
			p.errback += p.archNeuronNet[1][neuronOut].weight[neuronHid] *
				p.weightDelta[0][neuronOut]
		}
		p.weightDelta[1] = append(p.weightDelta[1], p.errback*
			p.archNeuronNet[0][neuronHid].output*
			(1-p.archNeuronNet[0][neuronHid].output))
		for neuronIn := 0; neuronIn <= len(input)-1; neuronIn++ {
			p.archNeuronNet[0][neuronHid].weight[neuronIn] = p.archNeuronNet[0][neuronHid].weight[neuronIn] -
				input[neuronIn]*
					p.learningRate*
					p.weightDelta[1][neuronHid]
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	PrintMemUsage()
	inputVals := []string{}
	var inputValsfloat []float64
	expectedVals := []string{}
	var expectedValsfloat []float64
	hiddenLayers := 5
	epochs := 1000

	i, err := ioutil.ReadFile("./education")
	check(err)
	flag := true
	countInpNeur := 0
	countExpNeur := 0
	for _, st := range i {
		if st == 44 {
			continue
		} else {
			if flag {
				if st == 9 {
					flag = false
					continue
				}
				countInpNeur++
			} else {
				if st == 10 {
					break
				}
				countExpNeur++
			}
		}
	}
	neuralNetwork := NeuralNetwork{
		learningRate: 0.9,
		numNeurons: []int{countInpNeur,
			hiddenLayers,
			countExpNeur,
		},
	}
	neuralNetwork.init()
	i1, err := os.Open("./test")
	check(err)
	defer i1.Close()
	r := bufio.NewReader(i1)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		flag = true
		a := scanner.Bytes()
		for _, val := range a {
			if val == 44 {
				continue
			} else {
				if flag {
					if val == 9 {
						flag = false
						continue
					}
					inputVals = append(inputVals, string(val))
				} else {
					if val != 10 {
						expectedVals = append(expectedVals, string(val))
					} else if val == 10 {
						break
					}
				}
			}
		}
		for _, val := range inputVals {
			value, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic("Error with parsing")
			}
			inputValsfloat = append(inputValsfloat, value)
		}

		for _, val := range expectedVals {
			value, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic("Error with parsing")
			}
			expectedValsfloat = append(expectedValsfloat, value)
		}

		neuralNetwork.dataStreamRight(inputValsfloat, expectedValsfloat, false)

		inputVals, inputValsfloat = inputVals[:0], inputValsfloat[:0]
		expectedVals, expectedValsfloat = expectedVals[:0], expectedValsfloat[:0]

	}
	fmt.Println("correctAnswer:", neuralNetwork.correctAnswer, "incorrectAnswer:", neuralNetwork.incorrectAnswer)

	neuralNetwork.correctAnswer = 0
	neuralNetwork.incorrectAnswer = 0
	PrintMemUsage()
	start := time.Now()
	for e := 0; e < epochs; e++ {
		i1, err := os.Open("./education")
		check(err)
		defer i1.Close()
		r := bufio.NewReader(i1)
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			flag = true
			a := scanner.Bytes()
			for _, val := range a {
				if val == 44 {
					continue
				} else {
					if flag {
						if val == 9 {
							flag = false
							continue
						}
						inputVals = append(inputVals, string(val))
					} else {
						if val != 10 {
							expectedVals = append(expectedVals, string(val))
						} else if val == 10 {
							break
						}
					}
				}
			}
			for _, val := range inputVals {
				value, err := strconv.ParseFloat(val, 64)
				if err != nil {
					panic("Error with parsing")
				}
				inputValsfloat = append(inputValsfloat, value)
			}

			for _, val := range expectedVals {
				value, err := strconv.ParseFloat(val, 64)
				if err != nil {
					panic("Error with parsing")
				}
				expectedValsfloat = append(expectedValsfloat, value)
			}

			neuralNetwork.dataStreamRight(inputValsfloat, expectedValsfloat, true)

			inputVals, inputValsfloat = inputVals[:0], inputValsfloat[:0]
			expectedVals, expectedValsfloat = expectedVals[:0], expectedValsfloat[:0]
			fmt.Println(time.Since(start))
		}

	}
	PrintMemUsage()
	i1, err = os.Open("./test")
	check(err)
	defer i1.Close()
	r = bufio.NewReader(i1)
	scanner = bufio.NewScanner(r)

	for scanner.Scan() {
		flag = true
		a := scanner.Bytes()
		for _, val := range a {
			if val == 44 {
				continue
			} else {
				if flag {
					if val == 9 {
						flag = false
						continue
					}
					inputVals = append(inputVals, string(val))
				} else {
					if val != 10 {
						expectedVals = append(expectedVals, string(val))
					} else if val == 10 {
						break
					}
				}
			}
		}
		for _, val := range inputVals {
			value, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic("Error with parsing")
			}
			inputValsfloat = append(inputValsfloat, value)
		}

		for _, val := range expectedVals {
			value, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic("Error with parsing")
			}
			expectedValsfloat = append(expectedValsfloat, value)
		}

		neuralNetwork.dataStreamRight(inputValsfloat, expectedValsfloat, false)

		inputVals, inputValsfloat = inputVals[:0], inputValsfloat[:0]
		expectedVals, expectedValsfloat = expectedVals[:0], expectedValsfloat[:0]

	}
	fmt.Println("correctAnswer:", neuralNetwork.correctAnswer, "incorrectAnswer:", neuralNetwork.incorrectAnswer)
	PrintMemUsage()
}

//PrintMemUsage using for tracking Memoty usage
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
