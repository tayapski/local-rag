package utils

func ConvertSlice(input []float64) []float32 {
    output := make([]float32, len(input))
    for i, v := range input {
        output[i] = float32(v)
    }
    return output
}
