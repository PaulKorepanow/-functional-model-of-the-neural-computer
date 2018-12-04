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

type Neuron struct {
	weight []float64
	input  []float64
	output float64
}

func random() float64 {
	weight := rand.Float64()
	return weight
}

func (n *Neuron) transfer_function(input []float64) {
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

type NeuralNetwork struct {
	learning_rate   float64
	numNeurons      []int
	archNeuronNet   [][]Neuron
	hiddenOut       []float64
	outputOut       []float64
	weight_delta    [][]float64
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
				val := random()
				ja.weight = append(ja.weight, val)
			}
			p.archNeuronNet[ia][j] = ja
		}
	}

	for ind := 0; ind <= neurons[1]; ind++ {
		p.err = make([]float64, neurons[1])
		val := random()
		p.err = append(p.err, val)
	}

}

type dataStream interface {
	data_stream_right(inputVals []int, expectedVals []int, b bool)
}

func (n *NeuralNetwork) data_stream_right(inputVals []float64, expectedVals []float64, edu bool) {
	for i := 0; i <= n.numNeurons[1]-1; i++ {
		n.archNeuronNet[0][i].transfer_function(inputVals)
		n.hiddenOut = append(n.hiddenOut, n.archNeuronNet[0][i].output)
	}
	for i := 0; i <= n.numNeurons[2]-1; i++ {
		n.archNeuronNet[1][i].transfer_function(n.hiddenOut)
		n.outputOut = append(n.outputOut, n.archNeuronNet[1][i].output)
	}
	n.hiddenOut = nil
	if edu {
		n.outputOut = nil
		n.data_stream_back(inputVals, expectedVals)
	} else {
		var compereVal1 float64
		var compereVal2 float64
		compereVal1 = 0
		compereVal2 = 0
		index1 := 0
		index2 := 0
		for i1, v1 := range expectedVals {
			if v1 > compereVal1 {
				compereVal1 = v1
				index2 = i1

			}
		}
		for i2, v2 := range n.outputOut {
			if v2 > compereVal2 {
				compereVal2 = v2
				index1 = i2
			}
		}
		if index1 == index2 {
			n.correctAnswer++
		} else {
			n.incorrectAnswer++
		}
	}
	n.outputOut = n.outputOut[:0]

}

func (n *NeuralNetwork) data_stream_back(input []float64, expectedVals []float64) {
	n.weight_delta = make([][]float64, 2)
	for neuronOut := 0; neuronOut <= n.numNeurons[2]-1; neuronOut++ {
		n.err[neuronOut] = n.archNeuronNet[1][neuronOut].output - expectedVals[neuronOut]
		delta := n.err[neuronOut] * n.archNeuronNet[1][neuronOut].output * (1 - n.archNeuronNet[1][neuronOut].output)
		n.weight_delta[0] = append(n.weight_delta[0], delta)

		for neuronHid := 0; neuronHid <= n.numNeurons[1]-1; neuronHid++ {
			n.archNeuronNet[1][neuronOut].weight[neuronHid] = n.archNeuronNet[1][neuronOut].weight[neuronHid] -
				n.archNeuronNet[0][neuronHid].output*
					n.learning_rate*
					delta

		}
	}
	for neuronHid := 0; neuronHid <= n.numNeurons[1]-1; neuronHid++ {
		n.errback = 0
		for neuronOut := 0; neuronOut <= n.numNeurons[2]-1; neuronOut++ {
			n.errback += n.archNeuronNet[1][neuronOut].weight[neuronHid] *
				n.weight_delta[0][neuronOut]
		}
		n.weight_delta[1] = append(n.weight_delta[1], n.errback*
			n.archNeuronNet[0][neuronHid].output*
			(1-n.archNeuronNet[0][neuronHid].output))
		for neuronIn := 0; neuronIn <= len(input)-1; neuronIn++ {
			n.archNeuronNet[0][neuronHid].weight[neuronIn] = n.archNeuronNet[0][neuronHid].weight[neuronIn] -
				input[neuronIn]*
					n.learning_rate*
					n.weight_delta[1][neuronHid]
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
} // тут доделать
//TODO тут доделать
func main() {
	PrintMemUsage()
	inputVals := []string{}
	var inputValsfloat []float64
	expectedVals := []string{}
	var expectedValsfloat []float64
	hiddenLayers := 10
	epochs := 1000
	var count int

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

	neural_network := NeuralNetwork{
		learning_rate: 0.9,
		numNeurons: []int{countInpNeur,
			hiddenLayers,
			countExpNeur,
		},
	}
	neural_network.init()
	i1, err := os.Open("./test")
	check(err)
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

		neural_network.data_stream_right(inputValsfloat, expectedValsfloat, false)

		inputVals, inputValsfloat = inputVals[:0], inputValsfloat[:0]
		expectedVals, expectedValsfloat = expectedVals[:0], expectedValsfloat[:0]

	}
	fmt.Println("correctAnswer:", neural_network.correctAnswer, "incorrectAnswer:", neural_network.incorrectAnswer)

	neural_network.correctAnswer = 0
	neural_network.incorrectAnswer = 0
	PrintMemUsage()
	start := time.Now()
	for e := 0; e < epochs; e++ {
		i1, err := os.Open("./education")
		check(err)
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
			count++
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

			neural_network.data_stream_right(inputValsfloat, expectedValsfloat, true)

			inputVals, inputValsfloat = inputVals[:0], inputValsfloat[:0]
			expectedVals, expectedValsfloat = expectedVals[:0], expectedValsfloat[:0]
			fmt.Println(time.Since(start))
		}

	}
	PrintMemUsage()
	i1, err = os.Open("./test")
	check(err)
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

		neural_network.data_stream_right(inputValsfloat, expectedValsfloat, false)

		inputVals, inputValsfloat = inputVals[:0], inputValsfloat[:0]
		expectedVals, expectedValsfloat = expectedVals[:0], expectedValsfloat[:0]

	}
	fmt.Println("correctAnswer:", neural_network.correctAnswer, "incorrectAnswer:", neural_network.incorrectAnswer)

}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
