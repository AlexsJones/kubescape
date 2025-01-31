package cautils

import (
	"path/filepath"

	"github.com/armosec/kubescape/cautils/getter"
	"github.com/armosec/opa-utils/reporthandling"
)

type ScanInfo struct {
	Getters
	PolicyIdentifier   []reporthandling.PolicyIdentifier
	UseExceptions      string   // Load file with exceptions configuration
	ControlsInputs     string   // Load file with inputs for controls
	UseFrom            []string // Load framework from local file (instead of download). Use when running offline
	UseDefault         bool     // Load framework from cached file (instead of download). Use when running offline
	Format             string   // Format results (table, json, junit ...)
	Output             string   // Store results in an output file, Output file name
	ExcludedNamespaces string   // DEPRECATED?
	IncludeNamespaces  string   // DEPRECATED?
	InputPatterns      []string // Yaml files input patterns
	Silent             bool     // Silent mode - Do not print progress logs
	FailThreshold      uint16   // Failure score threshold
	Submit             bool     // Submit results to Armo BE
	Local              bool     // Do not submit results
	Account            string   // account ID
	FrameworkScan      bool     // false if scanning control
	ScanAll            bool     // true if scan all frameworks
}

type Getters struct {
	ExceptionsGetter     getter.IExceptionsGetter
	ControlsInputsGetter getter.IControlsInputsGetter
	PolicyGetter         getter.IPolicyGetter
}

func (scanInfo *ScanInfo) Init() {
	scanInfo.setUseFrom()
	scanInfo.setUseExceptions()
	scanInfo.setAccountConfig()
	scanInfo.setOutputFile()

}

func (scanInfo *ScanInfo) setUseExceptions() {
	if scanInfo.UseExceptions != "" {
		// load exceptions from file
		scanInfo.ExceptionsGetter = getter.NewLoadPolicy([]string{scanInfo.UseExceptions})
	} else {
		scanInfo.ExceptionsGetter = getter.GetArmoAPIConnector()
	}
}

func (scanInfo *ScanInfo) setAccountConfig() {
	if scanInfo.ControlsInputs != "" {
		// load account config from file
		scanInfo.ControlsInputsGetter = getter.NewLoadPolicy([]string{scanInfo.ControlsInputs})
	} else {
		scanInfo.ControlsInputsGetter = getter.GetArmoAPIConnector()
	}
}
func (scanInfo *ScanInfo) setUseFrom() {
	if scanInfo.UseDefault {
		for _, policy := range scanInfo.PolicyIdentifier {
			scanInfo.UseFrom = append(scanInfo.UseFrom, getter.GetDefaultPath(policy.Name+".json"))
		}
	}
}

func (scanInfo *ScanInfo) setOutputFile() {
	if scanInfo.Output == "" {
		return
	}
	if scanInfo.Format == "json" {
		if filepath.Ext(scanInfo.Output) != ".json" {
			scanInfo.Output += ".json"
		}
	}
	if scanInfo.Format == "junit" {
		if filepath.Ext(scanInfo.Output) != ".xml" {
			scanInfo.Output += ".xml"
		}
	}
}

func (scanInfo *ScanInfo) ScanRunningCluster() bool {
	return len(scanInfo.InputPatterns) == 0
}

func (scanInfo *ScanInfo) SetPolicyIdentifiers(policies []string, kind reporthandling.NotificationPolicyKind) {
	for _, policy := range policies {
		if !scanInfo.contains(policy) {
			newPolicy := reporthandling.PolicyIdentifier{}
			newPolicy.Kind = kind // reporthandling.KindFramework
			newPolicy.Name = policy
			scanInfo.PolicyIdentifier = append(scanInfo.PolicyIdentifier, newPolicy)
		}
	}
}

func (scanInfo *ScanInfo) contains(policyName string) bool {
	for _, policy := range scanInfo.PolicyIdentifier {
		if policy.Name == policyName {
			return true
		}
	}
	return false
}
