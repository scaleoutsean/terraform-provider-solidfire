package elementsw

// buildScheduleInfo converts a slice of volume IDs to the correct scheduleInfo field for API requests
func buildScheduleInfo(volumes []string, retention string, snapMirrorLabel interface{}) map[string]interface{} {
	info := make(map[string]interface{})
	info["retention"] = retention
	if snapMirrorLabel != nil {
		info["snapMirrorLabel"] = snapMirrorLabel
	}
	if len(volumes) == 1 {
		info["volumeID"] = volumes[0]
	} else if len(volumes) > 1 {
		info["volumes"] = volumes
	}
	return info
}

// toIntSlice converts an interface{} list to []int
func toIntSlice(v interface{}) []int {
	if v == nil {
		return nil
	}
	arr, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]int, len(arr))
	for i, x := range arr {
		out[i] = x.(int)
	}
	return out
}
