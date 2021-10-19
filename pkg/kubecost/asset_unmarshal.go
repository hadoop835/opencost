package kubecost

import (
	"bytes"
	"fmt"
	"reflect"
	"time"

	// gojson is default golang json, required for RawMessage decoding
	gojson "encoding/json"

	"github.com/kubecost/cost-model/pkg/util/json"
)

// Encoding and decoding logic for Asset types

// Any marshal and unmarshal

// MarshalJSON implements json.Marshaler
func (a *Any) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncode(buffer, "properties", a.Properties(), ",")
	jsonEncode(buffer, "labels", a.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", a.AssetPricingModels(), ",")
	jsonEncode(buffer, "window", a.Window(), ",")
	jsonEncodeString(buffer, "start", a.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", a.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", a.Minutes(), ",")
	jsonEncodeFloat64(buffer, "adjustment", a.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", a.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", a.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", a.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (a *Any) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = a.InterfaceToAny(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to Any, carrying over relevant fields
func (a *Any) InterfaceToAny(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	a.properties = &properties
	a.labels = labels
	a.assetPricingModels = pricingModels
	a.start = start
	a.end = end
	a.window = Window{
		start: &start,
		end:   &end,
	}

	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		a.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		a.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		a.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		a.Cost = Cost.(float64) - a.adjustment - a.credit
	}

	return nil
}

// Cloud marshal and unmarshal

// MarshalJSON implements json.Marshaler
func (ca *Cloud) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", ca.Type().String(), ",")
	jsonEncode(buffer, "properties", ca.Properties(), ",")
	jsonEncode(buffer, "labels", ca.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", ca.AssetPricingModels(), ",")
	jsonEncode(buffer, "usageType", ca.UsageType, ",")
	jsonEncode(buffer, "usageDetail", ca.UsageDetail, ",")
	jsonEncode(buffer, "window", ca.Window(), ",")
	jsonEncodeString(buffer, "start", ca.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", ca.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", ca.Minutes(), ",")
	jsonEncodeFloat64(buffer, "adjustment", ca.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", ca.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", ca.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", ca.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (ca *Cloud) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = ca.InterfaceToCloud(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to Cloud, carrying over relevant fields
func (ca *Cloud) InterfaceToCloud(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	ca.properties = &properties
	ca.labels = labels
	ca.assetPricingModels = pricingModels
	ca.start = start
	ca.end = end
	ca.window = Window{
		start: &start,
		end:   &end,
	}

	if usageType, err := getTypedVal(fmap["usageType"]); err == nil {
		ca.UsageType = usageType.(string)
	}
	if usageDetail, err := getTypedVal(fmap["usageDetail"]); err == nil {
		ca.UsageDetail = usageDetail.(string)
	}
	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		ca.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		ca.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		ca.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		ca.Cost = Cost.(float64) - ca.adjustment - ca.credit
	}

	return nil
}

// ClusterManagement marshal and unmarshal

// MarshalJSON implements json.Marshler
func (cm *ClusterManagement) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", cm.Type().String(), ",")
	jsonEncode(buffer, "properties", cm.Properties(), ",")
	jsonEncode(buffer, "labels", cm.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", cm.AssetPricingModels(), ",")
	jsonEncode(buffer, "window", cm.Window(), ",")
	jsonEncodeString(buffer, "start", cm.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", cm.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", cm.Minutes(), ",")
	jsonEncodeFloat64(buffer, "adjustment", cm.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", cm.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", cm.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", cm.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (cm *ClusterManagement) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = cm.InterfaceToClusterManagement(f)
	if err != nil {
		return err
	}

	return nil
}

// InterfaceToClusterManagement Converts interface{} to ClusterManagement, carrying over relevant fields
func (cm *ClusterManagement) InterfaceToClusterManagement(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)


	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	cm.properties = &properties
	cm.labels = labels
	cm.assetPricingModels = pricingModels
	cm.window = Window{
		start: &start,
		end:   &end,
	}

	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		cm.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		cm.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		cm.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		cm.Cost = Cost.(float64) - cm.adjustment - cm.credit
	}

	return nil
}

// Disk marshal and unmarshal

// MarshalJSON implements the json.Marshaler interface
func (d *Disk) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", d.Type().String(), ",")
	jsonEncode(buffer, "properties", d.Properties(), ",")
	jsonEncode(buffer, "labels", d.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", d.AssetPricingModels(), ",")
	jsonEncode(buffer, "window", d.Window(), ",")
	jsonEncodeString(buffer, "start", d.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", d.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", d.Minutes(), ",")
	jsonEncodeFloat64(buffer, "byteHours", d.ByteHours, ",")
	jsonEncodeFloat64(buffer, "bytes", d.Bytes(), ",")
	jsonEncode(buffer, "breakdown", d.Breakdown, ",")
	jsonEncode(buffer, "storageClass", d.StorageClass, ",")
	jsonEncodeFloat64(buffer, "adjustment", d.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", d.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", d.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", d.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (d *Disk) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = d.InterfaceToDisk(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to Disk, carrying over relevant fields
func (d *Disk) InterfaceToDisk(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	fbreakdown := fmap["breakdown"].(map[string]interface{})

	breakdown := toBreakdown(fbreakdown)

	d.properties = &properties
	d.labels = labels
	d.assetPricingModels = pricingModels
	d.start = start
	d.end = end
	d.window = Window{
		start: &start,
		end:   &end,
	}
	d.Breakdown = &breakdown

	if storageClass, err := getTypedVal(fmap["storageClass"]); err == nil {
		d.StorageClass = storageClass.(string)
	}
	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		d.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		d.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		d.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		d.Cost = Cost.(float64) - d.adjustment - d.credit
	}
	if ByteHours, err := getTypedVal(fmap["byteHours"]); err == nil {
		d.ByteHours = ByteHours.(float64)
	}

	// d.Local is not marhsaled, and cannot be calculated from marshaled values.
	// Currently, it is just ignored and not set in the resulting unmarshal to Disk
	//  be aware that this means a resulting Disk from an unmarshal is therefore NOT
	// equal to the originally marshaled Disk.

	return nil

}

// Network marshal and unmarshal

// MarshalJSON implements json.Marshal interface
func (n *Network) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", n.Type().String(), ",")
	jsonEncode(buffer, "properties", n.Properties(), ",")
	jsonEncode(buffer, "labels", n.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", n.AssetPricingModels(), ",")
	jsonEncode(buffer, "window", n.Window(), ",")
	jsonEncodeString(buffer, "start", n.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", n.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", n.Minutes(), ",")
	jsonEncodeFloat64(buffer, "adjustment", n.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", n.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", n.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", n.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (n *Network) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = n.InterfaceToNetwork(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to Network, carrying over relevant fields
func (n *Network) InterfaceToNetwork(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	n.properties = &properties
	n.labels = labels
	n.assetPricingModels = pricingModels
	n.start = start
	n.end = end
	n.window = Window{
		start: &start,
		end:   &end,
	}

	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		n.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		n.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		n.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		n.Cost = Cost.(float64) - n.adjustment - n.credit
	}

	return nil

}

// Node marshal and unmarshal

// MarshalJSON implements json.Marshal interface
func (n *Node) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", n.Type().String(), ",")
	jsonEncode(buffer, "properties", n.Properties(), ",")
	jsonEncode(buffer, "assetPricingModels", n.AssetPricingModels(), ",")
	jsonEncode(buffer, "labels", n.Labels(), ",")
	jsonEncode(buffer, "window", n.Window(), ",")
	jsonEncodeString(buffer, "start", n.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", n.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", n.Minutes(), ",")
	jsonEncodeString(buffer, "nodeType", n.NodeType, ",")
	jsonEncodeFloat64(buffer, "cpuCores", n.CPUCores(), ",")
	jsonEncodeFloat64(buffer, "ramBytes", n.RAMBytes(), ",")
	jsonEncodeFloat64(buffer, "cpuCoreHours", n.CPUCoreHours, ",")
	jsonEncodeFloat64(buffer, "ramByteHours", n.RAMByteHours, ",")
	jsonEncodeFloat64(buffer, "GPUHours", n.GPUHours, ",")
	jsonEncode(buffer, "cpuBreakdown", n.CPUBreakdown, ",")
	jsonEncode(buffer, "ramBreakdown", n.RAMBreakdown, ",")
	jsonEncodeFloat64(buffer, "cpuCost", n.CPUCost, ",")
	jsonEncodeFloat64(buffer, "gpuCost", n.GPUCost, ",")
	jsonEncodeFloat64(buffer, "gpuCount", n.GPUs(), ",")
	jsonEncodeFloat64(buffer, "ramCost", n.RAMCost, ",")
	jsonEncodeFloat64(buffer, "adjustment", n.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", n.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", n.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", n.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (n *Node) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = n.InterfaceToNode(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to Node, carrying over relevant fields
func (n *Node) InterfaceToNode(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	fcpuBreakdown := fmap["cpuBreakdown"].(map[string]interface{})
	framBreakdown := fmap["ramBreakdown"].(map[string]interface{})

	cpuBreakdown := toBreakdown(fcpuBreakdown)
	ramBreakdown := toBreakdown(framBreakdown)

	n.properties = &properties
	n.labels = labels
	n.assetPricingModels = pricingModels
	n.start = start
	n.end = end
	n.window = Window{
		start: &start,
		end:   &end,
	}
	n.CPUBreakdown = &cpuBreakdown
	n.RAMBreakdown = &ramBreakdown

	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		n.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		n.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		n.discount = discount.(float64)
	}
	if NodeType, err := getTypedVal(fmap["nodeType"]); err == nil {
		n.NodeType = NodeType.(string)
	}
	if CPUCoreHours, err := getTypedVal(fmap["cpuCoreHours"]); err == nil {
		n.CPUCoreHours = CPUCoreHours.(float64)
	}
	if RAMByteHours, err := getTypedVal(fmap["ramByteHours"]); err == nil {
		n.RAMByteHours = RAMByteHours.(float64)
	}
	if GPUHours, err := getTypedVal(fmap["GPUHours"]); err == nil {
		n.GPUHours = GPUHours.(float64)
	}
	if CPUCost, err := getTypedVal(fmap["cpuCost"]); err == nil {
		n.CPUCost = CPUCost.(float64)
	}
	if GPUCost, err := getTypedVal(fmap["gpuCost"]); err == nil {
		n.GPUCost = GPUCost.(float64)
	}
	if GPUCount, err := getTypedVal(fmap["gpuCount"]); err == nil {
		n.GPUCount = GPUCount.(float64)
	}
	if RAMCost, err := getTypedVal(fmap["ramCost"]); err == nil {
		n.RAMCost = RAMCost.(float64)
	}

	return nil
}

// Loadbalancer marshal and unmarshal

// MarshalJSON implements json.Marshal
func (lb *LoadBalancer) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", lb.Type().String(), ",")
	jsonEncode(buffer, "properties", lb.Properties(), ",")
	jsonEncode(buffer, "labels", lb.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", lb.AssetPricingModels(), ",")
	jsonEncode(buffer, "window", lb.Window(), ",")
	jsonEncodeString(buffer, "start", lb.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", lb.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", lb.Minutes(), ",")
	jsonEncodeFloat64(buffer, "adjustment", lb.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", lb.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", lb.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", lb.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (lb *LoadBalancer) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = lb.InterfaceToLoadBalancer(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to LoadBalancer, carrying over relevant fields
func (lb *LoadBalancer) InterfaceToLoadBalancer(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	lb.properties = &properties
	lb.labels = labels
	lb.assetPricingModels = pricingModels
	lb.start = start
	lb.end = end
	lb.window = Window{
		start: &start,
		end:   &end,
	}

	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		lb.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		lb.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		lb.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		lb.Cost = Cost.(float64) - lb.adjustment - lb.credit
	}

	return nil

}

// SharedAsset marshal and unmarshal

// MarshalJSON implements json.Marshaler
func (sa *SharedAsset) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonEncodeString(buffer, "type", sa.Type().String(), ",")
	jsonEncode(buffer, "properties", sa.Properties(), ",")
	jsonEncode(buffer, "labels", sa.Labels(), ",")
	jsonEncode(buffer, "assetPricingModels", sa.AssetPricingModels(), ",")
	jsonEncode(buffer, "window", sa.Window(), ",")
	jsonEncodeString(buffer, "start", sa.Start().Format(time.RFC3339), ",")
	jsonEncodeString(buffer, "end", sa.End().Format(time.RFC3339), ",")
	jsonEncodeFloat64(buffer, "minutes", sa.Minutes(), ",")
	jsonEncodeFloat64(buffer, "adjustment", sa.Adjustment(), ",")
	jsonEncodeFloat64(buffer, "credit", sa.Credit(), ",")
	jsonEncodeFloat64(buffer, "discount", sa.Discount(), ",")
	jsonEncodeFloat64(buffer, "totalCost", sa.TotalCost(), "")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshal
func (sa *SharedAsset) UnmarshalJSON(b []byte) error {

	var f interface{}

	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	err = sa.InterfaceToSharedAsset(f)
	if err != nil {
		return err
	}

	return nil
}

// Converts interface{} to SharedAsset, carrying over relevant fields
func (sa *SharedAsset) InterfaceToSharedAsset(itf interface{}) error {

	fmap := itf.(map[string]interface{})

	// parse properties map to AssetProperties
	fproperties := fmap["properties"].(map[string]interface{})
	properties := toAssetProp(fproperties)

	// parse pricingModels map to AssetPricingModels
	var pricingModels *AssetPricingModels = nil
	if fpricingModels := fmap["assetPricingModels"]; fpricingModels != nil {
		pricingModels = toAssetPricingModels(fpricingModels.(map[string]interface{}))
	}

	// parse labels map to AssetLabels
	labels := make(map[string]string)
	for k, v := range fmap["labels"].(map[string]interface{}) {
		labels[k] = v.(string)
	}

	// parse start and end strings to time.Time
	start, err := time.Parse(time.RFC3339, fmap["start"].(string))
	if err != nil {
		return err
	}
	end, err := time.Parse(time.RFC3339, fmap["end"].(string))
	if err != nil {
		return err
	}

	sa.properties = &properties
	sa.labels = labels
	sa.assetPricingModels = pricingModels
	sa.window = Window{
		start: &start,
		end:   &end,
	}

	if adjustment, err := getTypedVal(fmap["adjustment"]); err == nil {
		sa.adjustment = adjustment.(float64)
	}
	if credit, err := getTypedVal(fmap["credit"]); err == nil {
		sa.credit = credit.(float64)
	}
	if discount, err := getTypedVal(fmap["discount"]); err == nil {
		sa.discount = discount.(float64)
	}
	if Cost, err := getTypedVal(fmap["totalCost"]); err == nil {
		sa.Cost = Cost.(float64) - sa.adjustment - sa.credit
	}

	return nil

}

// AssetSet marshal

// MarshalJSON JSON-encodes the AssetSet
func (as *AssetSet) MarshalJSON() ([]byte, error) {
	as.RLock()
	defer as.RUnlock()
	return json.Marshal(as.assets)
}

// AssetSetResponse for unmarshaling of AssetSet.assets into AssetSet

// UnmarshalJSON Unmarshals a marshaled AssetSet json into AssetSetResponse
func (asr *AssetSetResponse) UnmarshalJSON(b []byte) error {

	// gojson used here, as jsonitter UnmarshalJSON won't work with RawMessage
	var assetMap map[string]*gojson.RawMessage

	// Partial unmarshal to map of json RawMessage
	err := gojson.Unmarshal(b, &assetMap)
	if err != nil {
		return err
	}

	err = asr.RawMessageToAssetSetResponse(assetMap)
	if err != nil {
		return err
	}

	return nil
}

func (asr *AssetSetResponse) RawMessageToAssetSetResponse(assetMap map[string]*gojson.RawMessage) error {

	newAssetMap := make(map[string]Asset)

	// For each item in asset map, unmarshal to appropriate type
	for key, rawMessage := range assetMap {

		var f interface{}

		err := json.Unmarshal(*rawMessage, &f)
		if err != nil {
			return err
		}

		fmap := f.(map[string]interface{})

		switch t := fmap["type"]; t {
		case "Cloud":

			var ca Cloud
			err := ca.InterfaceToCloud(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &ca

		case "ClusterManagement":

			var cm ClusterManagement
			err := cm.InterfaceToClusterManagement(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &cm

		case "Disk":

			var d Disk
			err := d.InterfaceToDisk(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &d

		case "Network":

			var nw Network
			err := nw.InterfaceToNetwork(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &nw

		case "Node":

			var n Node
			err := n.InterfaceToNode(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &n

		case "LoadBalancer":

			var lb LoadBalancer
			err := lb.InterfaceToLoadBalancer(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &lb

		case "Shared":

			var sa SharedAsset
			err := sa.InterfaceToSharedAsset(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &sa

		default:

			var a Any
			err := a.InterfaceToAny(f)

			if err != nil {
				return err
			}

			newAssetMap[key] = &a

		}
	}

	asr.Assets = newAssetMap

	return nil
}

// UnmarshalJSON implements json.Unmarshal
func (asrr *AssetSetRangeResponse) UnmarshalJSON(b []byte) error {

	// gojson used here, as jsonitter UnmarshalJSON won't work with RawMessage
	var assetMapList []map[string]*gojson.RawMessage

	// Partial unmarshal to map of json RawMessage
	err := gojson.Unmarshal(b, &assetMapList)
	if err != nil {
		return err
	}

	var assetSetList []*AssetSetResponse

	for _, rawm := range assetMapList {

		var asresp AssetSetResponse
		err = asresp.RawMessageToAssetSetResponse(rawm)
		if err != nil {
			return err
		}

		assetSetList = append(assetSetList, &asresp)

	}

	asrr.Assets = assetSetList

	return nil
}

// Extra decoding util functions, for clarity

// Creates an AssetProperties directly from map[string]interface{}
func toAssetProp(fproperties map[string]interface{}) AssetProperties {
	var properties AssetProperties

	if category, v := fproperties["category"].(string); v {
		properties.Category = category
	}
	if provider, v := fproperties["provider"].(string); v {
		properties.Provider = provider
	}
	if account, v := fproperties["account"].(string); v {
		properties.Account = account
	}
	if project, v := fproperties["project"].(string); v {
		properties.Project = project
	}
	if service, v := fproperties["service"].(string); v {
		properties.Service = service
	}
	if cluster, v := fproperties["cluster"].(string); v {
		properties.Cluster = cluster
	}
	if name, v := fproperties["name"].(string); v {
		properties.Name = name
	}
	if providerID, v := fproperties["providerID"].(string); v {
		properties.ProviderID = providerID
	}
	if region, v := fproperties["region"].(string); v {
		properties.Region = region
	}
	if pricingSource, v := fproperties["pricingSource"].(string); v {
		properties.PricingSource = pricingSource
	}
	if currency, v := fproperties["currency"].(string); v {
		properties.Currency = currency
	}

	return properties
}

// Creates an AssetPricingModels directly from map[string]interface{}
func toAssetPricingModels(fpricingModels map[string]interface{}) *AssetPricingModels {
	if fpricingModels == nil {
		return nil
	}

	pricingModels := AssetPricingModels{}

	if preemptible, v := fpricingModels["preemptible"].(float64); v {
		pricingModels.Preemptible = preemptible
	}
	if reservedInstance, v := fpricingModels["reservedInstance"].(float64); v {
		pricingModels.ReservedInstance = reservedInstance
	}
	if savingsPlan, v := fpricingModels["savingsPlan"].(float64); v {
		pricingModels.SavingsPlan = savingsPlan
	}

	return &pricingModels
}

// Creates an Breakdown directly from map[string]interface{}
func toBreakdown(fproperties map[string]interface{}) Breakdown {
	var breakdown Breakdown

	if idle, v := fproperties["idle"].(float64); v {
		breakdown.Idle = idle
	}
	if other, v := fproperties["other"].(float64); v {
		breakdown.Other = other
	}
	if system, v := fproperties["system"].(float64); v {
		breakdown.System = system
	}
	if user, v := fproperties["user"].(float64); v {
		breakdown.User = user
	}

	return breakdown

}

// Not strictly nessesary, but cleans up the code and is a secondary check
// for correct types
func getTypedVal(itf interface{}) (interface{}, error) {
	switch itf := itf.(type) {
	case float64:
		return float64(itf), nil
	case string:
		return string(itf), nil
	default:
		unktype := reflect.ValueOf(itf)
		return nil, fmt.Errorf("Type %v is an invalid type", unktype)
	}
}
