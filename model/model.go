package model

type PCSpec struct {
	CPU         string  `json:"cpu"`         // e.g., "Intel Core i7-13700K"
	GPU         string  `json:"gpu"`         // e.g., "NVIDIA RTX 4070 Ti"
	RAM         int     `json:"ram"`         // RAM in GB, e.g., 32
	Storage     []Drive `json:"storage"`     // List of storage drives
	Motherboard string  `json:"motherboard"` // e.g., "ASUS ROG STRIX Z790-E"
	PSU         string  `json:"psu"`         // e.g., "Corsair RM850x 850W"
	Cooling     string  `json:"cooling"`     // e.g., "Noctua NH-D15"
	Case        string  `json:"case"`        // e.g., "Fractal Design Meshify 2"
	OS          string  `json:"os"`          // e.g., "Windows 11 Pro"
}

type Drive struct {
	Type     string `json:"type"`     // e.g., "SSD", "HDD"
	Capacity int    `json:"capacity"` // Capacity in GB
	Brand    string `json:"brand"`    // e.g., "Samsung", "Western Digital"
}

type Environment struct {
	Platform string `json:"platform"` // e.g., "Steam", "PS4", "RCPS3"
	Hardware string `json:"hardware"` // e.g., "PC", "PS4"
	PCSpec   PCSpec
}
