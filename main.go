package main

import (
    "fmt"
    "math"
    "math/rand"
)

const maxHideCount = 100

//BPNN BP                        
type BPNN struct {
    //            
    SampleCount int
    //                  
    InputCount int
    //                  
    OutputCount int
    //                                 
    HideCount int
    //         
    StudyRate float64
    //                  
    Precision float64
    //            
    LoopCount int
    //                  ,                           100
    hideWeight [][]float64
    //                  
    outWeight [][]float64
}

//NewBPNN             BP      
func NewBPNN(sc int, ic int, oc int, hc int, sr float64, p float64, lc int) (bp BPNN) {

    bp.SampleCount = sc
    bp.InputCount = ic
    bp.OutputCount = oc
    bp.HideCount = hc
    bp.StudyRate = sr
    bp.Precision = p
    bp.LoopCount = lc

    bp.hideWeight = make([][]float64, ic)
    bp.outWeight = make([][]float64, hc)

    //                     
    var i, j int
    for i = 0; i < ic; i++ {
        bp.hideWeight[i] = make([]float64, hc)
        for j = 0; j < hc; j++ {
            bp.hideWeight[i][j] = rand.Float64()
        }
    }

    for i = 0; i < hc; i++ {
        bp.outWeight[i] = make([]float64, oc)
        for j = 0; j < oc; j++ {
            bp.outWeight[i][j] = rand.Float64()
        }
    }

    return
}

//TrainBP       bp                  x                  y
func (bp *BPNN) TrainBP(x [][]float64, y [][]int) {
    //                  
    prec := bp.Precision
    //         
    studyRate := bp.StudyRate
    //               
    hideCount := bp.HideCount
    //                  
    loopCount := bp.LoopCount
    //                  
    hideWeight := bp.hideWeight
    //                  
    outWeight := bp.outWeight

    //               
    var ChgH = make([]float64, hideCount)
    var ChgO = make([]float64, bp.OutputCount)

    //                           
    var O1 = make([]float64, hideCount)
    var O2 = make([]float64, bp.OutputCount)

    //            
    var i, j, m, n int

    e := prec + 1
    //                                                      
    for n = 0; e > prec && n < loopCount; n++ {
        e = 0
        //                           
        for i = 0; i < bp.SampleCount; i++ {
            //                        
            for m = 0; m < hideCount; m++ {
                temp := 0.0
                for j = 0; j < bp.InputCount; j++ {
                    //fmt.Println(i, j, m, n, temp)
                    temp = temp + x[i][j]*hideWeight[j][m]
                }
                O1[m] = sigmoid(temp)
            }
            //                           
            for m = 0; m < bp.OutputCount; m++ {
                temp := 0.0
                for j = 0; j < hideCount; j++ {
                    temp = temp + O1[j]*outWeight[j][m]
                }
                O2[m] = sigmoid(temp)
            }
            
            //                              
            for j = 0; j < bp.OutputCount; j++ {
                ChgO[j] = O2[j] * (1 - O2[j]) * (float64(y[i][j]) - O2[j])
            }
            
            //                  
            for j = 0; j < bp.OutputCount; j++ {
                e = e + (float64(y[i][j])-O2[j])*(float64(y[i][j])-O2[j])
            }
            
            //                        
            for j = 0; j < hideCount; j++ {
                temp := 0.0
                for m = 0; m < bp.OutputCount; m++ {
                    temp = temp + outWeight[j][m]*ChgO[m]
                }
                ChgH[j] = temp * O1[j] * (1 - O1[j])
            }
            
            //                        
            for j = 0; j < hideCount; j++ {
                for m = 0; m < bp.OutputCount; m++ {
                    outWeight[j][m] = outWeight[j][m] + studyRate*O1[j]*ChgO[m]
                }
            }
            //                     
            for j = 0; j < bp.InputCount; j++ {
                for m = 0; m < hideCount; m++ {
                    hideWeight[j][m] = hideWeight[j][m] + studyRate*x[i][j]*ChgH[m]
                }
            }
        }
        //                                 
        if n%1000 == 0 {
            fmt.Printf("   %d                   : %f\n", n, e)
        }
    }
    
    fmt.Println("                     %d", n)
    fmt.Println("                              ")
    for i = 0; i < bp.InputCount; i++ {
        for j = 0; j < hideCount; j++ {
            fmt.Printf("%f\t", hideWeight[i][j])
        }
        fmt.Println(" ")
    }
    fmt.Println("                                 ")
    for i = 0; i < hideCount; i++ {
        for j = 0; j < bp.OutputCount; j++ {
            fmt.Printf("%f\t", outWeight[i][j])
        }
        fmt.Println(" ")
    }

    fmt.Println("BP                  !")
}

//Sigmoid      ,                        
func sigmoid(net float64) float64 {
    return 1 / (1 + math.Exp(-net))
}

//UseBP       bp                  
func (bp *BPNN) UseBP(input []float64) (output []float64) {
    var O1 = make([]float64, maxHideCount)
    var i, j int
    var temp float64

    output = make([]float64, bp.OutputCount)

    for i = 0; i < bp.HideCount; i++ {
        temp = 0
        for j = 0; j < bp.InputCount; j++ {
            temp = temp + input[j]*bp.hideWeight[j][i]
        }
        O1[i] = sigmoid(temp)
    }
    for i = 0; i < bp.OutputCount; i++ {
        temp = 0
        for j = 0; j < bp.HideCount; j++ {
            temp = temp + O1[j]*bp.outWeight[j][i]
        }
        output[i] = sigmoid(temp)
    }
    fmt.Println("         ")
    for i = 0; i < bp.OutputCount; i++ {
        fmt.Printf("%f\t", output[i])
    }
    fmt.Println(" ")
    return
}

func main() {
    //            
    x := [][]float64{
        []float64{0.8, 0.5, 0},
        []float64{0.9, 0.7, 0.3},
        []float64{1, 0.8, 0.5},
        []float64{0, 0.2, 0.3},
        []float64{0.2, 0.1, 1.3},
        []float64{0.2, 0.7, 0.8},
    }

    //            
    y := [][]int{
        []int{0, 1},
        []int{0, 1},
        []int{0, 1},
        []int{1, 0},
        []int{1, 0},
        []int{1, 0},
    }
    
    //            
    sampleCount := len(x)
    //                  
    inputCount := len(x[0])
    //                  
    outputCount := len(y[0])
    //                                 
    hideCount := 10
    //         
    studyRate := 0.01
    //                  
    precision := 0.001
    //            
    loopCount := 1000000
        
    bp := NewBPNN(sampleCount, inputCount, outputCount, hideCount, studyRate, precision, loopCount)
    bp.TrainBP(x, y)
    
    input := []float64{0.8, 0.8, 0}
    bp.UseBP(input)
}