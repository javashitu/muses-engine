package bo

var (
	P4k   Resolution = Resolution{Name: "4k", Height: 2160, Alias: "原画"}
	P1080 Resolution = Resolution{Name: "1080p", Height: 1080, Alias: "超清"}
	P720  Resolution = Resolution{Name: "720", Height: 720, Alias: "高清"}
	P480  Resolution = Resolution{Name: "480", Height: 480, Alias: "清晰"}
	P360  Resolution = Resolution{Name: "360", Height: 360, Alias: "流畅"}
)

type Resolution struct {
	Name   string
	Height int
	Alias  string
}

func WrapResolution(height int) Resolution {
	if height < 480 {
		return P360
	} else if height < 720 {
		return P480
	} else if height < 1080 {
		return P720
	} else if height < 2160 {
		return P1080
	} else {
		return P4k
	}
}
