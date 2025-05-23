package activedirectory

import (
	"github.com/lkarlslund/adalanche/modules/engine"
)

var (
	EdgeForeignIdentity = engine.NewEdge("ForeignIdentity")

	DistinguishedName                       = engine.NewAttribute("distinguishedName").Tag("AD").Unique().Single()
	ObjectClass                             = engine.NewAttribute("objectClass").Tag("AD")
	ObjectCategory                          = engine.NewAttribute("objectCategory").Tag("AD").Single()
	Type                                    = engine.NewAttribute("type").Single()
	StructuralObjectClass                   = engine.NewAttribute("structuralObjectClass").Tag("AD")
	NTSecurityDescriptor                    = engine.NewAttribute("nTSecurityDescriptor").Tag("AD").Single().Type(engine.AttributeTypeSecurityDescriptor)
	SAMAccountType                          = engine.NewAttribute("sAMAccountType").Tag("AD").Single()
	GroupType                               = engine.NewAttribute("groupType").Tag("AD").Single()
	MemberOf                                = engine.NewAttribute("memberOf").Tag("AD")
	Member                                  = engine.NewAttribute("member").Tag("AD")
	BadPasswordTime                         = engine.NewAttribute("badPasswordTime").Tag("AD").Type(engine.AttributeTypeTime100NS)
	MsPKIRoamingTimeStamp                   = engine.NewAttribute("msPKIRoamingTimeStamp").Tag("AD").Type(engine.AttributeTypeTime100NS)
	CreationTime                            = engine.NewAttribute("creationTime").Tag("AD").Type(engine.AttributeTypeTime100NS)
	AccountExpires                          = engine.NewAttribute("accountExpires").Tag("AD").Type(engine.AttributeTypeTime100NS)
	RepsTo                                  = engine.NewAttribute("repsTo").Tag("AD")
	InstanceType                            = engine.NewAttribute("instanceType").Tag("AD")
	ModifiedCount                           = engine.NewAttribute("modifiedCount").Tag("AD")
	MinPwdAge                               = engine.NewAttribute("minPwdAge").Tag("AD")
	MinPwdLength                            = engine.NewAttribute("minPwdLength").Tag("AD").Type(engine.AttributeTypeInt)
	PwdProperties                           = engine.NewAttribute("pwdProperties").Tag("AD")
	LockOutDuration                         = engine.NewAttribute("lockoutDuration").Tag("AD")
	PwdHistoryLength                        = engine.NewAttribute("pwdHistoryLength").Tag("AD").Type(engine.AttributeTypeInt)
	IsCriticalSystemObject                  = engine.NewAttribute("isCriticalSystemObject").Tag("AD")
	FSMORoleOwner                           = engine.NewAttribute("fSMORoleOwner").Tag("AD")
	NTMixedDomain                           = engine.NewAttribute("nTMixedDomain").Tag("AD")
	SystemFlags                             = engine.NewAttribute("systemFlags").Tag("AD")
	PrimaryGroupID                          = engine.NewAttribute("primaryGroupID").Tag("AD")
	LogonCount                              = engine.NewAttribute("logonCount").Tag("AD")
	UserAccountControl                      = engine.NewAttribute("userAccountControl").Tag("AD")
	LocalPolicyFlags                        = engine.NewAttribute("localPolicyFlags").Tag("AD")
	CodePage                                = engine.NewAttribute("codePage").Tag("AD")
	CountryCode                             = engine.NewAttribute("countryCode").Tag("AD")
	OperatingSystem                         = engine.NewAttribute("operatingSystem").Tag("AD")
	OperatingSystemHotfix                   = engine.NewAttribute("operatingSystemHotfix").Tag("AD")
	OperatingSystemVersion                  = engine.NewAttribute("operatingSystemVersion").Tag("AD")
	OperatingSystemServicePack              = engine.NewAttribute("operatingSystemServicePack").Tag("AD")
	AdminCount                              = engine.NewAttribute("adminCount").Tag("AD")
	LogonHours                              = engine.NewAttribute("logonHours").Tag("AD")
	BadPwdCount                             = engine.NewAttribute("badPwdCount").Tag("AD").Type(engine.AttributeTypeInt)
	GPCFileSysPath                          = engine.NewAttribute("gPCFileSysPath").Tag("AD").Merge()
	SchemaIDGUID                            = engine.NewAttribute("schemaIDGUID").Tag("AD").Type(engine.AttributeTypeGUID)
	PossSuperiors                           = engine.NewAttribute("possSuperiors")
	SystemPossSuperiors                     = engine.NewAttribute("systemPossSuperiors")
	SubClassOf                              = engine.NewAttribute("subClassOf").Tag("AD")
	SystemMayContain                        = engine.NewAttribute("systemMayContain")
	SystemMustContain                       = engine.NewAttribute("systemMustContain")
	ServicePrincipalName                    = engine.NewAttribute("servicePrincipalName").Tag("AD")
	Name                                    = engine.NewAttribute("name").Tag("AD")
	DisplayName                             = engine.NewAttribute("displayName").Tag("AD")
	LDAPDisplayName                         = engine.NewAttribute("lDAPDisplayName").Tag("AD") // Attribute-Schema
	Description                             = engine.NewAttribute("description").Tag("AD")
	SAMAccountName                          = engine.NewAttribute("sAMAccountName").Tag("AD")
	ObjectSid                               = engine.NewAttribute("objectSid").Tag("AD").Merge().Single().Type(engine.AttributeTypeSID)
	CreatorSID                              = engine.NewAttribute("mS-DS-CreatorSID").Tag("AD").Single().Type(engine.AttributeTypeSID)
	MSDSAllowedToActOnBehalfOfOtherIdentity = engine.NewAttribute("msDS-AllowedToActOnBehalfOfOtherIdentity").Tag("AD").Type(engine.AttributeTypeSecurityDescriptor)
	MSDSAllowedToDelegateTo                 = engine.NewAttribute("msDS-AllowedToDelegateTo").Tag("AD")
	FRSRootSecurity                         = engine.NewAttribute("fRSRootSecurity").Tag("AD").Type(engine.AttributeTypeSecurityDescriptor)
	MSDFSLinkSecurityDescriptorv2           = engine.NewAttribute("msDFS-LinkSecurityDescriptorv2").Tag("AD").Type(engine.AttributeTypeSecurityDescriptor)
	PKIEnrollmentAccess                     = engine.NewAttribute("pKIEnrollmentAccess").Tag("AD")
	ObjectGUID                              = engine.NewAttribute("objectGUID").Tag("AD").Merge()
	PwdLastSet                              = engine.NewAttribute("pwdLastSet").Tag("AD").Type(engine.AttributeTypeTime)
	WhenCreated                             = engine.NewAttribute("whenCreated").Type(engine.AttributeTypeTime)
	WhenChanged                             = engine.NewAttribute("whenChanged").Type(engine.AttributeTypeTime)
	DsCorePropagationData                   = engine.NewAttribute("dsCorePropagationData").Type(engine.AttributeTypeTime)
	MsExchLastUpdateTime                    = engine.NewAttribute("msExchLastUpdateTime").Type(engine.AttributeTypeTime)
	GWARTLastModified                       = engine.NewAttribute("gWARTLastModified").Type(engine.AttributeTypeTime)
	SpaceLastComputed                       = engine.NewAttribute("spaceLastComputed").Type(engine.AttributeTypeTime)
	MsExchPolicyLastAppliedTime             = engine.NewAttribute("msExchPolicyLastAppliedTime").Type(engine.AttributeTypeTime)
	MsExchWhenMailboxCreated                = engine.NewAttribute("msExchWhenMailboxCreated").Type(engine.AttributeTypeTime)
	SIDHistory                              = engine.NewAttribute("sIDHistory").Tag("AD").Type(engine.AttributeTypeSID)
	LastLogon                               = engine.NewAttribute("lastLogon").Type(engine.AttributeTypeTime)
	LastLogonTimestamp                      = engine.NewAttribute("lastLogonTimestamp").Type(engine.AttributeTypeTime)
	MSDSGroupMSAMembership                  = engine.NewAttribute("msDS-GroupMSAMembership").Tag("AD").Type(engine.AttributeTypeSecurityDescriptor)
	MSDSHostServiceAccount                  = engine.NewAttribute("msDS-HostServiceAccount").Tag("AD")
	MSDSHostServiceAccountBL                = engine.NewAttribute("msDS-HostServiceAccountBL").Tag("AD")

	MSmcsAdmPwdExpirationTime    = engine.NewAttribute("ms-mcs-AdmPwdExpirationTime").Tag("AD").Type(engine.AttributeTypeTime)    // LAPS V1 password timeout
	MSLAPSPasswordExpirationTime = engine.NewAttribute("ms-LAPS-PasswordExpirationTime").Tag("AD").Type(engine.AttributeTypeTime) // LAPS V2 password timeout

	SecurityIdentifier       = engine.NewAttribute("securityIdentifier").Type(engine.AttributeTypeSID)
	TrustDirection           = engine.NewAttribute("trustDirection").Type(engine.AttributeTypeInt)
	TrustAttributes          = engine.NewAttribute("trustAttributes")
	TrustPartner             = engine.NewAttribute("trustPartner")
	TrustType                = engine.NewAttribute("trustType")
	DsHeuristics             = engine.NewAttribute("dsHeuristics").Tag("AD")
	AttributeSecurityGUID    = engine.NewAttribute("attributeSecurityGUID").Tag("AD")
	MSDSConsistencyGUID      = engine.NewAttribute("mS-DS-ConsistencyGuid")
	RightsGUID               = engine.NewAttribute("rightsGUID").Tag("AD").Type(engine.AttributeTypeGUID)
	GPLink                   = engine.NewAttribute("gPLink").Tag("AD")
	GPOptions                = engine.NewAttribute("gPOptions").Tag("AD")
	ScriptPath               = engine.NewAttribute("scriptPath").Tag("AD").Single()
	MSPKICertificateNameFlag = engine.NewAttribute("msPKI-Certificate-Name-Flag").Tag("AD").Type(engine.AttributeTypeInt)
	PKIExtendedUsage         = engine.NewAttribute("pKIExtendedKeyUsage").Tag("AD")
	PKIExpirationPeriod      = engine.NewAttribute("pKIExpirationPeriod").Tag("AD")
	PKIOverlapPeriod         = engine.NewAttribute("pKIOverlapPeriod").Tag("AD")
	MsDSBehaviourVersion     = engine.NewAttribute("msDS-Behavior-Version").Type(engine.AttributeTypeInt)
	DNSHostName              = engine.NewAttribute("dnsHostName").Tag("AD")

	MetaPasswordAge  = engine.NewAttribute("passwordAge").SetDescription("Password age in hours")
	MetaLastLoginAge = engine.NewAttribute("lastLoginAge").SetDescription("Time elapsed since last login in hours")
)
