package activedirectory

import "github.com/lkarlslund/adalanche/modules/engine"

func OnlyIfTargetAccountEnabled(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
	if target.HasTag("account_enabled") || edges.IsSet(EdgeWriteUserAccountControl) {
		return 100
	}
	return 0
}

func FixedProbability(probability int) engine.ProbabilityCalculatorFunction {
	return func(source, target *engine.Object, edge *engine.EdgeBitmap) engine.Probability {
		return engine.Probability(probability)
	}
}

var (
	EdgeACLContainsDeny  = engine.NewEdge("ACLContainsDeny").RegisterProbabilityCalculator(FixedProbability(0)).Tag("Informative")
	EdgeResetPassword    = engine.NewEdge("ResetPassword").RegisterProbabilityCalculator(OnlyIfTargetAccountEnabled).Tag("Pivot")
	EdgeReadPasswordId   = engine.NewEdge("ReadPasswordId").SetDefault(false, false, false).RegisterProbabilityCalculator(FixedProbability(5))
	EdgeOwns             = engine.NewEdge("Owns").Tag("Pivot")
	EdgeGenericAll       = engine.NewEdge("GenericAll").Tag("Informative")
	EdgeWriteAll         = engine.NewEdge("WriteAll").Tag("Informative").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeWritePropertyAll = engine.NewEdge("WritePropertyAll").Tag("Informative").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeWriteExtendedAll = engine.NewEdge("WriteExtendedAll").Tag("Informative").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeTakeOwnership    = engine.NewEdge("TakeOwnership").Tag("Pivot")
	EdgeWriteDACL        = engine.NewEdge("WriteDACL").Tag("Pivot")
	EdgeWriteSPN         = engine.NewEdge("WriteSPN").RegisterProbabilityCalculator(func(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
		if target.HasTag("account_active") {
			return 50
		}
		return 0
	}).Tag("Pivot")
	EdgeWriteValidatedSPN = engine.NewEdge("WriteValidatedSPN").RegisterProbabilityCalculator(func(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
		if target.HasTag("account_active") {
			return 50
		}
		return 0
	}).Tag("Pivot")
	EdgeWriteAllowedToAct        = engine.NewEdge("WriteAllowedToAct").Tag("Pivot")
	EdgeWriteAllowedToDelegateTo = engine.NewEdge("WriteAllowedToDelegTo").Tag("Pivot")
	EdgeAddMember                = engine.NewEdge("AddMember").Tag("Pivot")
	EdgeAddMemberGroupAttr       = engine.NewEdge("AddMemberGroupAttr").Tag("Pivot")
	EdgeAddSelfMember            = engine.NewEdge("AddSelfMember").Tag("Pivot")
	EdgeReadGMSAPassword         = engine.NewEdge("ReadGMSAPassword").Tag("Pivot")
	EdgeHasMSA                   = engine.NewEdge("HasMSA").Tag("Granted")
	EdgeWriteUserAccountControl  = engine.NewEdge("WriteUserAccountControl").Describe("Allows attacker to set ENABLE and set DONT_REQ_PREAUTH and then to do AS_REP Kerberoasting").RegisterProbabilityCalculator(func(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
		/*if uac, ok := target.AttrInt(activedirectory.UserAccountControl); ok && uac&0x0002 != 0 { //UAC_ACCOUNTDISABLE
			// Account is disabled
			return 0
		}*/
		return 50
	}).Tag("Pivot")

	EdgeWriteKeyCredentialLink = engine.NewEdge("WriteKeyCredentialLink").RegisterProbabilityCalculator(func(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
		if target.HasTag("account_enabled") || edges.IsSet(EdgeWriteUserAccountControl) {
			return 100
		}
		return 0
	}).Tag("Pivot")
	EdgeWriteAttributeSecurityGUID           = engine.NewEdge("WriteAttrSecurityGUID").RegisterProbabilityCalculator(FixedProbability(0)) // Only if you patch the DC, so this will actually never work
	EdgeSIDHistoryEquality                   = engine.NewEdge("SIDHistoryEquality").Tag("Pivot")
	EdgeAllExtendedRights                    = engine.NewEdge("AllExtendedRights").Tag("Informative").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeDSReplicationSyncronize              = engine.NewEdge("DSReplSync").Tag("Granted").SetDefault(false, false, false).Tag("Granted").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeDSReplicationGetChanges              = engine.NewEdge("DSReplGetChngs").SetDefault(false, false, false).Tag("Granted").Tag("Granted").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeDSReplicationGetChangesAll           = engine.NewEdge("DSReplGetChngsAll").SetDefault(false, false, false).Tag("Granted").Tag("Granted").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeDSReplicationGetChangesInFilteredSet = engine.NewEdge("DSReplGetChngsInFiltSet").SetDefault(false, false, false).Tag("Granted").Tag("Granted").RegisterProbabilityCalculator(FixedProbability(0))
	EdgeCall                                 = engine.NewEdge("Call").Describe("Call a service point")
	EdgeControls                             = engine.NewEdge("Controls").Describe("Node controls a service point")
	EdgeReadLAPSPassword                     = engine.NewEdge("ReadLAPSPassword").Tag("Pivot").Tag("Granted")
	EdgeMemberOfGroup                        = engine.NewEdge("MemberOfGroup").Tag("Granted")
	EdgeMemberOfGroupIndirect                = engine.NewEdge("MemberOfGroupIndirect").SetDefault(false, false, false).Tag("Granted")
	EdgeHasSPN                               = engine.NewEdge("HasSPN").Describe("Kerberoastable by requesting Kerberos service ticket against SPN and then bruteforcing the ticket").RegisterProbabilityCalculator(func(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
		if target.HasTag("account_enabled") || edges.IsSet(EdgeWriteUserAccountControl) {
			return 50
		}
		// Account is disabled
		return 0
	}).Tag("Pivot")
	EdgeDontReqPreauth = engine.NewEdge("DontReqPreauth").Describe("Kerberoastable by AS-REP by requesting a TGT and then bruteforcing the ticket").RegisterProbabilityCalculator(func(source, target *engine.Object, edges *engine.EdgeBitmap) engine.Probability {
		if target.HasTag("account_enabled") || edges.IsSet(EdgeWriteUserAccountControl) {
			return 50
		}
		return 0
	}).Tag("Pivot")
	EdgeOverwritesACL              = engine.NewEdge("OverwritesACL")
	EdgeAffectedByGPO              = engine.NewEdge("AffectedByGPO").Tag("Granted").Tag("Pivot")
	PartOfGPO                      = engine.NewEdge("PartOfGPO").Tag("Granted").Tag("Pivot")
	EdgeLocalAdminRights           = engine.NewEdge("AdminRights").Tag("Granted").Tag("Pivot")
	EdgeLocalRDPRights             = engine.NewEdge("RDPRights").RegisterProbabilityCalculator(FixedProbability(30)).Tag("Pivot")
	EdgeLocalDCOMRights            = engine.NewEdge("DCOMRights").RegisterProbabilityCalculator(FixedProbability(30)).Tag("Pivot")
	EdgeScheduledTaskOnUNCPath     = engine.NewEdge("SchedTaskOnUNCPath").Tag("Pivot")
	EdgeMachineScript              = engine.NewEdge("MachineScript").Tag("Pivot")
	EdgeWriteAltSecurityIdentities = engine.NewEdge("WriteAltSecIdent").Tag("Pivot").RegisterProbabilityCalculator(OnlyIfTargetAccountEnabled)
	EdgeWriteProfilePath           = engine.NewEdge("WriteProfilePath").Tag("Pivot")
	EdgeWriteScriptPath            = engine.NewEdge("WriteScriptPath").Tag("Pivot")
	EdgeCertificateEnroll          = engine.NewEdge("CertificateEnroll").Tag("Granted")
	EdgeCertificateAutoEnroll      = engine.NewEdge("CertificateAutoEnroll").Tag("Granted")
	EdgeVoodooBit                  = engine.NewEdge("VoodooBit").SetDefault(false, false, false).Tag("Internal").Hidden()
)
