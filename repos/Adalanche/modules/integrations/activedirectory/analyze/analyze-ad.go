package analyze

import (
	"encoding/binary"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gofrs/uuid"
	"github.com/lkarlslund/adalanche/modules/engine"
	"github.com/lkarlslund/adalanche/modules/integrations/activedirectory"
	"github.com/lkarlslund/adalanche/modules/ui"
	"github.com/lkarlslund/adalanche/modules/util"
	"github.com/lkarlslund/adalanche/modules/windowssecurity"
	"github.com/puzpuzpuz/xsync"
)

// Interesting permissions on AD
var (
	ResetPwd, _                             = uuid.FromString("{00299570-246d-11d0-a768-00aa006e0529}")
	DSReplicationGetChanges                 = uuid.UUID{0x11, 0x31, 0xf6, 0xaa, 0x9c, 0x07, 0x11, 0xd1, 0xf7, 0x9f, 0x00, 0xc0, 0x4f, 0xc2, 0xdc, 0xd2}
	DSReplicationGetChangesAll              = uuid.UUID{0x11, 0x31, 0xf6, 0xad, 0x9c, 0x07, 0x11, 0xd1, 0xf7, 0x9f, 0x00, 0xc0, 0x4f, 0xc2, 0xdc, 0xd2}
	DSReplicationSyncronize                 = uuid.UUID{0x11, 0x31, 0xf6, 0xab, 0x9c, 0x07, 0x11, 0xd1, 0xf7, 0x9f, 0x00, 0xc0, 0x4f, 0xc2, 0xdc, 0xd2}
	DSReplicationGetChangesInFilteredSet, _ = uuid.FromString("{89e95b76-444d-4c62-991a-0facbeda640c}")

	AttributeMember                = uuid.UUID{0xbf, 0x96, 0x79, 0xc0, 0x0d, 0xe6, 0x11, 0xd0, 0xa2, 0x85, 0x00, 0xaa, 0x00, 0x30, 0x49, 0xe2}
	AttributeSetGroupMembership, _ = uuid.FromString("{BC0AC240-79A9-11D0-9020-00C04FC2D4CF}")
	AttributeSIDHistory            = uuid.UUID{0x17, 0xeb, 0x42, 0x78, 0xd1, 0x67, 0x11, 0xd0, 0xb0, 0x02, 0x00, 0x00, 0xf8, 0x03, 0x67, 0xc1}

	AttributeAllowedToActOnBehalfOfOtherIdentity, _ = uuid.FromString("{3F78C3E5-F79A-46BD-A0B8-9D18116DDC79}")
	AttributeAllowedToDelegateTo, _                 = uuid.FromString("{800d94d7-b7a1-42a1-b14d-7cae1423d07f}")

	AttributeMSDSGroupMSAMembership       = uuid.UUID{0x88, 0x8e, 0xed, 0xd6, 0xce, 0x04, 0xdf, 0x40, 0xb4, 0x62, 0xb8, 0xa5, 0x0e, 0x41, 0xba, 0x38}
	AttributeGPLink, _                    = uuid.FromString("{F30E3BBE-9FF0-11D1-B603-0000F80367C1}")
	AttributeMSDSKeyCredentialLink, _     = uuid.FromString("{5B47D60F-6090-40B2-9F37-2A4DE88F3063}")
	AttributeSecurityGUIDGUID, _          = uuid.FromString("{bf967924-0de6-11d0-a285-00aa003049e2}")
	AttributeAltSecurityIdentitiesGUID, _ = uuid.FromString("{00FBF30C-91FE-11D1-AEBC-0000F80367C1}")
	AttributeProfilePathGUID, _           = uuid.FromString("{bf967a05-0de6-11d0-a285-00aa003049e2}")
	AttributeScriptPathGUID, _            = uuid.FromString("{bf9679a8-0de6-11d0-a285-00aa003049e2}")
	AttributeMSDSManagedPasswordId, _     = uuid.FromString("{0e78295a-c6d3-0a40-b491-d62251ffa0a6}")
	AttributeUserAccountControlGUID, _    = uuid.FromString("{bf967a68-0de6-11d0-a285-00aa003049e2}")
	AttributePwdLastSetGUID, _            = uuid.FromString("{bf967a0a-0de6-11d0-a285-00aa003049e2}")

	ExtendedRightCertificateEnroll, _     = uuid.FromString("{0e10c968-78fb-11d2-90d4-00c04f79dc55}")
	ExtendedRightCertificateAutoEnroll, _ = uuid.FromString("{a05b8cc2-17bc-4802-a710-e7c15ab866a2}")

	ValidateWriteSelfMembership, _ = uuid.FromString("{bf9679c0-0de6-11d0-a285-00aa003049e2}")
	ValidateWriteSPN, _            = uuid.FromString("{f3a64788-5306-11d1-a9c5-0000f80367c1}")

	ObjectGuidUser, _            = uuid.FromString("{bf967aba-0de6-11d0-a285-00aa003049e2")
	ObjectGuidComputer, _        = uuid.FromString("{bf967a86-0de6-11d0-a285-00aa003049e2")
	ObjectGuidGroup, _           = uuid.FromString("{bf967a9c-0de6-11d0-a285-00aa003049e2")
	ObjectGuidDomain, _          = uuid.FromString("{19195a5a-6da0-11d0-afd3-00c04fd930c9")
	ObjectGuidDNSZone, _         = uuid.FromString("{e0fa1e8b-9b45-11d0-afdd-00c04fd930c9")
	ObjectGuidDNSNode, _         = uuid.FromString("{e0fa1e8c-9b45-11d0-afdd-00c04fd930c9")
	ObjectGuidGPO, _             = uuid.FromString("{f30e3bc2-9ff0-11d1-b603-0000f80367c1")
	ObjectGuidOU, _              = uuid.FromString("{bf967aa5-0de6-11d0-a285-00aa003049e2")
	ObjectGuidAttributeSchema, _ = uuid.FromString("{BF967A80-0DE6-11D0-A285-00AA003049E2}")

	AdministratorsSID, _           = windowssecurity.ParseStringSID("S-1-5-32-544")
	BackupOperatorsSID, _          = windowssecurity.ParseStringSID("S-1-5-32-551")
	PrintOperatorsSID, _           = windowssecurity.ParseStringSID("S-1-5-32-550")
	ServerOperatorsSID, _          = windowssecurity.ParseStringSID("S-1-5-32-549")
	EnterpriseDomainControllers, _ = windowssecurity.ParseStringSID("S-1-5-9")

	GPLinkCache = engine.NewAttribute("gpLinkCache")

	NetBIOSName = engine.NewAttribute("nETBIOSName")
	NCName      = engine.NewAttribute("nCName")
	DNSRoot     = engine.NewAttribute("dnsRoot")

	MemberOfIndirect = engine.NewAttribute("memberOfIndirect")

	ObjectTypeMachine    = engine.NewObjectType("Machine", "Machine")
	DomainJoinedSID      = engine.NewAttribute("domainJoinedSid").Merge()
	DnsHostName          = engine.NewAttribute("dnsHostName")
	EdgeAuthenticatesAs  = engine.NewEdge("AuthenticatesAs")
	EdgeInheritsSecurity = engine.NewEdge("InheritsSecurity").SetDefault(true, true, false)

	CertificateTemplates   = engine.NewAttribute("certificateTemplates")
	PublishedBy            = engine.NewAttribute("publishedBy")
	PublishedByDnsHostName = engine.NewAttribute("publishedByDnsHostName")

	MetaPasswordAge  = engine.NewAttribute("passwordAge")
	MetaLastLoginAge = engine.NewAttribute("lastLoginAge")

	msLAPSEncryptedPasswordAttributesGUID, _ = uuid.FromString("{f3531ec6-6330-4f8e-8d39-7a671fbac605}")

	EdgeMachineAccount = engine.NewEdge("MachineAccount").RegisterProbabilityCalculator(activedirectory.FixedProbability(-1)).Describe("Indicates this is the domain joined computer account belonging to the machine")
)

var warnedgpos = make(map[string]struct{})

func init() {
	engine.AddMergeApprover("Only merge Machine objects with other Machine objects", func(a, b *engine.Object) (*engine.Object, error) {
		if a.Type() == ObjectTypeMachine && b.Type() != ObjectTypeMachine {
			return nil, engine.ErrDontMerge
		} else if b.Type() == ObjectTypeMachine && a.Type() != ObjectTypeMachine {
			return nil, engine.ErrDontMerge
		}
		return nil, nil
	})

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Find LAPS or return
		var lapsGUID uuid.UUID
		if lapsobject, found := ao.FindTwo(engine.Name, engine.NewAttributeValueString("ms-Mcs-AdmPwd"),
			engine.ObjectClass, engine.NewAttributeValueString("attributeSchema")); found {
			if objectGUID, ok := lapsobject.OneAttrRaw(activedirectory.SchemaIDGUID).(uuid.UUID); ok {
				ui.Debug().Msg("Detected LAPS schema extension GUID")
				lapsGUID = objectGUID
			} else {
				ui.Error().Msgf("Could not read LAPS schema extension GUID from %v", lapsobject.DN())
			}
		}

		if lapsGUID.IsNil() {
			ui.Debug().Msg("Microsoft LAPS V1 not detected, skipping tests for this")
			return
		}

		ao.Iterate(func(o *engine.Object) bool {
			// Only for computers
			if o.Type() != engine.ObjectTypeComputer {
				return true
			}

			// ... that has LAPS installed
			if !o.HasAttr(activedirectory.MSmcsAdmPwdExpirationTime) {
				return true
			}

			// Analyze ACL
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}

			// Link to the machine object
			machinesid := o.SID()
			if machinesid.IsBlank() {
				ui.Fatal().Msgf("Computer account %v has no objectSID", o.DN())
			}
			machine, found := ao.Find(DomainJoinedSID, engine.NewAttributeValueSID(machinesid))
			if !found {
				ui.Error().Msgf("Could not locate machine for domain SID %v while processing LAPS v1", machinesid)
				return true
			}
			machine.Tag("laps")

			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, lapsGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(machine, activedirectory.EdgeReadLAPSPassword)
				}
			}
			return true
		})
	}, "Reading local admin passwords via LAPS v1", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Find LAPS or return
		var lapsV2PasswordGUID uuid.UUID
		var lapsV2EncryptedPasswordGUID uuid.UUID

		if lapsobject, found := ao.FindTwo(engine.Name, engine.NewAttributeValueString("ms-LAPS-Password"),
			engine.ObjectClass, engine.NewAttributeValueString("attributeSchema")); found {
			if objectGUID, ok := lapsobject.OneAttrRaw(activedirectory.SchemaIDGUID).(uuid.UUID); ok {
				ui.Debug().Msg("Detected LAPS schema extension GUID")
				lapsV2PasswordGUID = objectGUID
			} else {
				ui.Error().Msgf("Could not read LAPS schema extension GUID from %v", lapsobject.DN())
			}
		}
		if lapsobject, found := ao.FindTwo(engine.Name, engine.NewAttributeValueString("ms-LAPS-EncryptedPassword"),
			engine.ObjectClass, engine.NewAttributeValueString("attributeSchema")); found {
			if objectGUID, ok := lapsobject.OneAttrRaw(activedirectory.SchemaIDGUID).(uuid.UUID); ok {
				ui.Debug().Msg("Detected LAPS schema extension GUID")
				lapsV2EncryptedPasswordGUID = objectGUID
			} else {
				ui.Error().Msgf("Could not read LAPS schema extension GUID from %v", lapsobject.DN())
			}
		}

		if lapsV2PasswordGUID.IsNil() {
			ui.Debug().Msg("Microsoft LAPS V2 not detected, skipping tests for this")
			return
		}

		ao.Iterate(func(o *engine.Object) bool {
			// Only for computers
			if o.Type() != engine.ObjectTypeComputer {
				return true
			}

			// ... that has LAPS installed
			if !o.HasAttr(activedirectory.MSLAPSPasswordExpirationTime) {
				return true
			}

			// Analyze ACL
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}

			// Link to the machine object
			machinesid := o.SID()
			if machinesid.IsBlank() {
				ui.Fatal().Msgf("Computer account %v has no objectSID", o.DN())
			}
			machine, found := ao.Find(DomainJoinedSID, engine.NewAttributeValueSID(machinesid))
			if !found {
				ui.Error().Msgf("Could not locate machine for domain SID %v while processing LAPS v2", machinesid)
				return true
			}
			machine.Tag("laps")

			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, lapsV2PasswordGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(machine, activedirectory.EdgeReadLAPSPassword)
				}
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, lapsV2EncryptedPasswordGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(machine, activedirectory.EdgeReadLAPSPassword) // FIXME
				}
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, msLAPSEncryptedPasswordAttributesGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(machine, activedirectory.EdgeReadLAPSPassword) // FIXME
				}
			}
			return true
		})
	}, "Reading local admin passwords via LAPS v2", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.Type() == engine.ObjectTypeForeignSecurityPrincipal {
				return true
			}
			if sd, err := o.SecurityDescriptor(); err == nil && sd.Control&engine.CONTROLFLAG_DACL_PROTECTED == 0 {
				if parentobject, found := ao.DistinguishedParent(o); found {
					parentobject.EdgeTo(o, EdgeInheritsSecurity)
				}
			}
			return true
		})
	}, "Indicator that object inherits security from the container it is within", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.Type() != engine.ObjectTypeContainer || o.OneAttrString(engine.Name) != "Machine" {
				return true
			}
			// Only for computers, you can't really pwn users this way
			p, hasparent := ao.DistinguishedParent(o)
			if !hasparent || p.Type() != engine.ObjectTypeGroupPolicyContainer {
				return true
			}
			p.EdgeTo(o, activedirectory.PartOfGPO)
			return true
		})
	}, "Machine configurations that are part of a GPO", engine.BeforeMergeHigh)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.Type() != engine.ObjectTypeContainer || o.OneAttrString(engine.Name) != "User" {
				return true
			}
			// Only for users, you can't really pwn users this way
			p, hasparent := ao.DistinguishedParent(o)
			if !hasparent || p.Type() != engine.ObjectTypeGroupPolicyContainer {
				return true
			}
			p.EdgeTo(o, activedirectory.PartOfGPO)
			return true
		})
	}, "User configurations that are part of a GPO", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for _, acl := range sd.DACL.Entries {
				if acl.Type == engine.ACETYPE_ACCESS_DENIED || acl.Type == engine.ACETYPE_ACCESS_DENIED_OBJECT {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeACLContainsDeny) // Not a probability of success, this is just an indicator
				}
			}
			return true
		})
	}, "Indicator for possible false positives, as the ACL contains DENY entries", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			// https://www.alsid.com/crb_article/kerberos-delegation/
			// --- Citation bloc --- This is generally true, but an exception exists: positioning a Deny for the OWNER RIGHTS SID (S-1-3-4) in an object’s ACE removes the owner’s implicit control of this object’s DACL. ---------------------
			aclhasdeny := false
			for _, ace := range sd.DACL.Entries {
				if ace.Type == engine.ACETYPE_ACCESS_DENIED && ace.SID == windowssecurity.OwnerSID {
					aclhasdeny = true
				}
			}
			if !sd.Owner.IsNull() && !aclhasdeny {
				ao.FindOrAddAdjacentSID(sd.Owner, o).EdgeTo(o, activedirectory.EdgeOwns)
			}
			return true
		})
	}, "Indicator that someone owns an object", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_GENERIC_ALL, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeGenericAll)
				}
			}
			return true
		})
	}, "Indicator that someone has full permissions on an object", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_GENERIC_WRITE, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteAll)
				}
			}
			return true
		})
	}, "Indicator that someone can write to all attributes and do all validated writes on an object", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWritePropertyAll)
				}
			}
			return true
		})
	}, "Indicator that someone can write to all attributes of an object", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY_EXTENDED, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteExtendedAll)
				}
			}
			return true
		})
	}, "Indicator that someone do all validated writes on an object", engine.BeforeMergeFinal)

	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-dtyp/c79a383c-2b3f-4655-abe7-dcbb7ce0cfbe IMPORTANT
	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_WRITE_OWNER, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeTakeOwnership)
				}
			}
			return true
		})
	}, "Indicator that someone is allowed to take ownership of an object", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_WRITE_DACL, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteDACL)
				}
			}
			return true
		})
	}, "Indicator that someone can change permissions on an object", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if o.Type() != engine.ObjectTypeAttributeSchema {
				return true
			}
			// FIXME - check for SYSTEM ATTRIBUTES - these can NEVER be changed
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeSecurityGUIDGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteAttributeSecurityGUID) // Experimental, I've never run into this misconfiguration
				}
			}
			return true
		})
	}, `Allows an attacker to modify the attribute security set of an attribute, promoting it to a weaker attribute set (experimental/wrong)`, engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only users, computers and service accounts
			if o.Type() != engine.ObjectTypeUser && o.Type() != engine.ObjectTypeComputer {
				return true
			}
			// Check who can reset the password
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, ResetPwd, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeResetPassword)
				}
			}
			return true
		})
	}, "Indicator that a group or user can reset the password of an account", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only group managed service accounts
			if o.Type() != engine.ObjectTypeGroupManagedServiceAccount {
				return true
			}

			// Check who can reset the password
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_READ_PROPERTY, AttributeMSDSManagedPasswordId, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeReadPasswordId)
				}
			}
			return true
		})
	}, "Indicator that a group or user can read the msDS-ManagedPasswordId for use in MGSA Golden attack", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		kerberoast := "kerberoast"
		authusers := FindWellKnown(ao, windowssecurity.AuthenticatedUsersSID)
		if authusers == nil {
			ui.Error().Msgf("Could not locate Authenticated Users")
			return
		}

		ao.Iterate(func(o *engine.Object) bool {
			// Only computers and users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			if o.Attr(activedirectory.ServicePrincipalName).Len() > 0 {
				o.Tag(kerberoast)
				authusers.EdgeTo(o, activedirectory.EdgeHasSPN)
			}
			return true
		})
	}, "Indicator that a user has a ServicePrincipalName and an authenticated user can Kerberoast it", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		anonymous := FindWellKnown(ao, windowssecurity.AnonymousLogonSID)

		if anonymous == nil {
			ui.Error().Msgf("Could not locate Anonymous Logon")
			return
		}

		ao.Iterate(func(o *engine.Object) bool {
			// Only users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			if uac, ok := o.AttrInt(activedirectory.UserAccountControl); ok && uac&engine.UAC_DONT_REQ_PREAUTH != 0 {
				o.Tag("asreproast")
				anonymous.EdgeTo(o, activedirectory.EdgeDontReqPreauth)
			}
			return true
		})
	}, "Indicator that a user has \"don't require preauth\" and can be ASREPRoasted", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, ValidateWriteSPN, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteSPN)
				}
			}
			return true
		})
	}, "Indicator that a user can change the ServicePrincipalName attribute, and then Kerberoast the account", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only computers and users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY_EXTENDED, ValidateWriteSPN, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteValidatedSPN)
				}
			}
			return true
		})
	}, "Indicator that a user can change the ServicePrincipalName attribute (validate write), and then Kerberoast the account", engine.BeforeMergeFinal)

	// https://blog.harmj0y.net/activedirectory/the-most-dangerous-user-right-you-probably-have-never-heard-of/
	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only computers
			if o.Type() != engine.ObjectTypeComputer && o.Type() != engine.ObjectTypeUser {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeAllowedToActOnBehalfOfOtherIdentity, ao) {
					// This does NOT requires the SeEnableDelegationPrivilege set on the DC for the user doing it!!
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteAllowedToAct)
				}
			}
			return true
		})
	}, `Modify the msDS-AllowedToActOnBehalfOfOtherIdentity (Resource Based Constrained Delegation) on an account to enable any SPN enabled user to impersonate it`, engine.BeforeMergeFinal)

	EdgeRBCD := engine.NewEdge("RBConstrainedDeleg")
	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only computers
			if o.Type() != engine.ObjectTypeComputer && o.Type() != engine.ObjectTypeUser {
				return true
			}
			o.Attr(activedirectory.MSDSAllowedToActOnBehalfOfOtherIdentity).Iterate(func(val engine.AttributeValue) bool {
				// Each of these is a SID, so find that SID and add an edge
				if sd, ok := val.Raw().(*engine.SecurityDescriptor); ok {
					// ui.Debug().Msgf("Found msDS-AllowedToActOnBehalfOfOtherIdentity on %v as %v", o.DN(), sd.String(ao))
					for _, acl := range sd.DACL.Entries {
						if acl.Type == engine.ACETYPE_ACCESS_ALLOWED {
							ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, EdgeRBCD)
						}
					}
				}
				return true
			})
			return true
		})
	}, `Someone is listed in the msDS-AllowedToActOnBehalfOfOtherIdentity (Resource Based Constrained Delegation) on an account`, engine.BeforeMergeFinal)

	EdgeCD := engine.NewEdge("ConstrainedDeleg")
	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only computers
			if o.Type() != engine.ObjectTypeComputer && o.Type() != engine.ObjectTypeUser {
				return true
			}
			if uac, ok := o.AttrInt(activedirectory.UserAccountControl); ok {
				if uac&engine.UAC_TRUSTED_TO_AUTH_FOR_DELEGATION != 0 {
					o.Attr(activedirectory.MSDSAllowedToDelegateTo).Iterate(func(val engine.AttributeValue) bool {
						// Each of these is a SID, so find that SID and add an edge
						// sd := val.Raw().(*engine.SecurityDescriptor)
						ui.Debug().Msgf("Found msDS-AllowedToDelegate on %v as %v", o.DN(), val.String())
						_, host, split := strings.Cut(val.String(), "/")
						if !split {
							ui.Error().Msgf("Constrained delegation SPN %v does not contain /", val.String())
							return true // continue
						}
						if strings.Contains(host, "/") {
							ui.Error().Msgf("Constrained delegation host name %v still contains /", val.String())
							return true // continue
						}
						if strings.Contains(host, ":") {
							ui.Debug().Msgf("Constrained delegation host name %v contains :, removing port", val.String())
							host = strings.Split(host, ":")[0]
						}
						if !strings.Contains(host, ".") {
							ui.Debug().Msgf("Constrained delegation host name %v is not FQDN, adding domain context DNS", val.String())
							host += "." + util.DomainContextToDomainSuffix(o.OneAttrString(engine.DomainContext))
						}
						if target, found := ao.FindTwo(DnsHostName, engine.NewAttributeValueString(host),
							engine.Type, engine.NewAttributeValueString("Machine"),
						); found {
							o.EdgeTo(target, EdgeCD)
						} else {
							ui.Error().Msgf("Could not find constrained delegation SPN %v target (looked for machine %v) in the AD", val.String(), host)
						}

						return true
					})
				}
			}
			return true
		})
	}, `Someone is listed in the msDS-AllowedToDelegate (Constrained Delegation) on an account`, engine.BeforeMergeFinal)

	/*
		// https://blog.harmj0y.net/activedirectory/the-most-dangerous-user-right-you-probably-have-never-heard-of/
		Loader.AddProcessor(func(ao *engine.Objects) {
			ao.Iterate(func(o *engine.Object) bool {
				// Only computers
				if o.Type() != engine.ObjectTypeComputer {
					return true
				}
				sd, err := o.SecurityDescriptor()
				if err != nil {
					return true
				}
				for index, acl := range sd.DACL.Entries {
					if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeAllowedToDelegateTo, ao) {
						// Also requires the SeEnableDelegationPrivilege set on the DC for the user doing it!!
						ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteAllowedToDelegateTo) // Success rate?
					}
				}
				return true
			})
		}, `Modify the msDS-AllowedToDelegateTo (Constrained Delegation) on a computer to enable any SPN enabled user to impersonate anyone else`, engine.BeforeMergeFinal)
	*/
	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for groups
			if o.Type() != engine.ObjectTypeGroup {
				return true
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeMember, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeAddMember)
				}
			}
			return true
		})
	}, "Permission to add a member to a group", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for groups
			if o.Type() != engine.ObjectTypeGroup {
				return true
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeSetGroupMembership, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeAddMemberGroupAttr)
				}
			}
			return true
		})
	}, "Permission to add a member to a group (via attribute set)", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for groups
			if o.Type() != engine.ObjectTypeGroup {
				return true
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY_EXTENDED, ValidateWriteSelfMembership, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeAddSelfMember)
				}
			}
			return true
		})
	}, "Permission to add yourself to a group", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			o.Attr(activedirectory.MSDSGroupMSAMembership).Iterate(func(msads engine.AttributeValue) bool {
				if sd, ok := msads.Raw().(*engine.SecurityDescriptor); ok {
					for _, acl := range sd.DACL.Entries {
						if acl.Type == engine.ACETYPE_ACCESS_ALLOWED {
							ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeReadGMSAPassword)
						}
					}
				}
				return true
			})
			return true
		})
	}, "Allows someone to read a password of a managed service account", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeAltSecurityIdentitiesGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteAltSecurityIdentities)
				}
			}
			return true
		})
	}, "Allows an attacker to define a certificate that can be used to authenticate as the user", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeProfilePathGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteProfilePath)
				}
			}
			return true
		})
	}, "Change user profile path (allows an attacker to trigger a user auth against an attacker controlled UNC path)", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for users
			if o.Type() != engine.ObjectTypeUser {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeScriptPathGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteScriptPath)
				}
			}
			return true
		})
	}, "Change user script path (allows an attacker to trigger a user auth against an attacker controlled UNC path)", engine.BeforeMergeFinal)
	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			o.Attr(activedirectory.MSDSHostServiceAccount).Iterate(func(dn engine.AttributeValue) bool {
				if targetmsa, found := ao.Find(engine.DistinguishedName, dn); found {
					o.EdgeTo(targetmsa, activedirectory.EdgeHasMSA)
				}
				return true
			})
			return true
		})
	}, "Indicates that the object has a service account in use", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for groups
			if o.Type() != engine.ObjectTypeUser && o.Type() != engine.ObjectTypeComputer {
				return true
			}
			// It's a group
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeMSDSKeyCredentialLink, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteKeyCredentialLink)
				}
			}
			return true
		})
	}, "Allows you to write your own cert to keyCredentialLink, and then auth as that user (no password reset needed)", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			o.Attr(activedirectory.SIDHistory).Iterate(func(sidval engine.AttributeValue) bool {
				if sid, ok := sidval.Raw().(windowssecurity.SID); ok {
					target := ao.FindOrAddAdjacentSID(sid, o)
					o.EdgeTo(target, activedirectory.EdgeSIDHistoryEquality)
				}
				return true
			})
			return true
		})
	}, "Indicates that object has a SID History attribute pointing to the other object, making them the 'same' permission wise", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeAllExtendedRights)
				}
			}
			return true
		})
	}, "Indicates that you have all extended rights", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.Type() != engine.ObjectTypeCertificateTemplate {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-crtd/211ab1e3-bad6-416d-9d56-8480b42617a4
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, ExtendedRightCertificateEnroll, ao) ||
					sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_VOODOO_BIT, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeCertificateEnroll)
				}
			}
			return true
		})
	}, "Permission to enroll into a certificate template", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.Type() != engine.ObjectTypeCertificateTemplate {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				// https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-crtd/211ab1e3-bad6-416d-9d56-8480b42617a4
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, ExtendedRightCertificateAutoEnroll, ao) ||
					sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_VOODOO_BIT, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeCertificateAutoEnroll)
				}
			}
			return true
		})
	}, "Permission to auto-enroll into a certificate template", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_VOODOO_BIT, uuid.Nil, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeVoodooBit)
				}
			}
			return true
		})
	}, "Has the Voodoo Bit set", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.Type() != engine.ObjectTypeDomainDNS {
				return true
			}
			if !o.HasAttr(activedirectory.SystemFlags) {
				return true
			}
			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}

			DCsyncObject, _ := ao.FindTwoOrAdd(
				engine.Type, engine.ObjectTypeCallableServicePoint.ValueString(),
				engine.Name, engine.NewAttributeValueString("DCsync"),
			)
			DCsyncObject.Tag("hvt")

			o.EdgeTo(DCsyncObject, activedirectory.EdgeControls)

			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, DSReplicationSyncronize, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeDSReplicationSyncronize)
				}
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, DSReplicationGetChanges, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeDSReplicationGetChanges)
				}
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, DSReplicationGetChangesAll, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeDSReplicationGetChangesAll)
				}
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_CONTROL_ACCESS, DSReplicationGetChangesInFilteredSet, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeDSReplicationGetChangesInFilteredSet)
				}
			}

			// Add the DCsync combination flag
			o.Edges(engine.In).Range(func(target *engine.Object, edge engine.EdgeBitmap) bool {
				if edge.IsSet(activedirectory.EdgeDSReplicationGetChanges) && edge.IsSet(activedirectory.EdgeDSReplicationGetChangesAll) {
					// DCsync attack WOT WOT
					target.EdgeTo(DCsyncObject, activedirectory.EdgeCall)
				}
				return true
			})
			return true
		})
	}, "Permissions on DomainDNS objects leading to DCsync attacks", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Ensure everyone has a family
		ao.Iterate(func(computeraccount *engine.Object) bool {
			if computeraccount.Type() != engine.ObjectTypeComputer {
				return true
			}

			sid := computeraccount.OneAttr(engine.ObjectSid)
			if sid == nil {
				ui.Error().Msgf("Computer account without SID: %v", computeraccount.DN())
				return true
			}
			machine, _ := ao.FindOrAdd(
				DomainJoinedSID, sid,
				engine.IgnoreBlanks,
				engine.Name, computeraccount.Attr(engine.Name),
				activedirectory.Type, ObjectTypeMachine.ValueString(),
				DnsHostName, computeraccount.Attr(DnsHostName),
			)
			// ui.Debug().Msgf("Added machine for SID %v", sid.String())

			machine.EdgeTo(computeraccount, EdgeAuthenticatesAs)
			machine.EdgeTo(computeraccount, EdgeMachineAccount)
			machine.ChildOf(computeraccount)

			return true
		})
	},
		"creating Machine objects (representing the machine running the OS)",
		engine.BeforeMerge)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Ensure everyone has a family
		ao.Iterate(func(o *engine.Object) bool {
			potentialorphan := o
			for {
				if potentialorphan.Parent() != nil {
					return true
				}

				if potentialorphan == ao.Root() {
					return true
				}

				if parent, found := ao.DistinguishedParent(potentialorphan); found {
					potentialorphan.ChildOf(parent)
					return true
				}

				dn := potentialorphan.DN()
				if potentialorphan.Type() == engine.ObjectTypeDomainDNS && len(dn) > 3 && strings.EqualFold("dc=", dn[:3]) {
					// Top of some AD we think, hook to top of browsable tree
					o.ChildOf(ao.Root())
					return true
				}

				// Create a synthetic parent
				parentdn := util.ParentDistinguishedName(potentialorphan.DN())
				if parentdn == "" {
					return true
				}

				ui.Debug().Msgf("AD object %v (%v) has no parent :-( - creating synthetic object", o.Label(), o.DN())

				newparent := ao.AddNew(
					engine.DistinguishedName, parentdn,
					engine.Description, "Synthetic parent object",
				)
				potentialorphan.ChildOf(newparent)
				potentialorphan = newparent // loop, to ensure new objects also have parents
			}
		})
	},
		"applying parent/child relationships",
		engine.BeforeMergeHigh)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		type domaininfo struct {
			suffix string
			name   string
		}
		var domains []domaininfo

		results, found := ao.FindMulti(engine.ObjectClass, engine.NewAttributeValueString("crossRef"))

		if !found {
			ui.Error().Msg("No domainDNS object found, can't apply DownLevelLogonName to objects")
			return
		}

		results.Iterate(func(o *engine.Object) bool {
			// Store domain -> netbios name in array for later
			dn := o.OneAttrString(NCName)
			netbiosname := o.OneAttrString(NetBIOSName)

			if dn == "" || netbiosname == "" {
				// Some crossref objects have no NCName or NetBIOSName, skip them
				return true // continue
			}

			domains = append(domains, domaininfo{
				suffix: dn,
				name:   netbiosname,
			})
			return true
		})

		if len(domains) == 0 {
			ui.Error().Msg("No NCName to NetBIOSName mapping found, can't apply DownLevelLogonName to objects")
			return
		}

		// Sort the domains so we match on longest first
		sort.Slice(domains, func(i, j int) bool {
			// Less is More - so we sort in reverse order
			return len(domains[i].suffix) > len(domains[j].suffix)
		})

		// Apply DownLevelLogonName to relevant objects
		ao.Iterate(func(o *engine.Object) bool {
			if !o.HasAttr(engine.SAMAccountName) {
				return true
			}
			dn := o.DN()
			for _, domaininfo := range domains {
				if strings.HasSuffix(dn, domaininfo.suffix) {
					o.Set(engine.DownLevelLogonName, engine.NewAttributeValueString(domaininfo.name+"\\"+o.OneAttrString(engine.SAMAccountName)))
					break
				}
			}
			return true
		})

		ao.DropIndex(engine.DownLevelLogonName)
	},
		"applying DownLevelLogonName attribute",
		engine.BeforeMergeLow)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Add domain part attribute from distinguished name to objects
		ao.Iterate(func(o *engine.Object) bool {
			// Only objects with a DistinguishedName
			if o.DN() == "" {
				return true
			}

			if o.HasAttr(engine.DomainContext) {
				return true
			}

			parts := strings.Split(o.DN(), ",")
			lastpart := -1

			for i := len(parts) - 1; i >= 0; i-- {
				part := parts[i]
				if len(part) < 3 || !strings.EqualFold("dc=", part[:3]) {
					break
				}
				if strings.EqualFold("DC=ForestDNSZones", part) || strings.EqualFold("DC=DomainDNSZones", part) {
					break
				}
				lastpart = i
			}

			if lastpart != -1 {
				o.Set(engine.DomainContext, engine.NewAttributeValueString(strings.Join(parts[lastpart:], ",")))
			}
			return true
		})
	},
		"applying domain part attribute",
		engine.BeforeMergeLow)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Find all the AdminSDHolder containers
		ao.Filter(func(o *engine.Object) bool {
			return strings.HasPrefix(o.OneAttrString(engine.DistinguishedName), "CN=AdminSDHolder,CN=System,")
		}).Iterate(func(adminsdholder *engine.Object) bool {
			// We found it - so we know it can change ACLs of some objects
			domaincontext := adminsdholder.OneAttrString(engine.DomainContext)

			// Are some groups excluded?
			excluded_mask := 0

			// Find dsHeuristics, this defines groups EXCLUDED From AdminSDHolder application
			// https://social.technet.microsoft.com/wiki/contents/articles/22331.adminsdholder-protected-groups-and-security-descriptor-propagator.aspx#What_is_a_protected_group
			if ds, found := ao.Find(engine.DistinguishedName, engine.NewAttributeValueString("CN=Directory Service,CN=Windows NT,CN=Services,CN=Configuration,"+domaincontext)); found {
				excluded := ds.OneAttrString(activedirectory.DsHeuristics)
				if len(excluded) >= 16 {
					excluded_mask = strings.Index("0123456789ABCDEF", strings.ToUpper(string(excluded[15])))
				}
			}

			ao.Filter(func(o *engine.Object) bool {
				// Check if object is a group
				if o.Type() != engine.ObjectTypeGroup {
					return false
				}

				// Only this "local" AD (for multi domain analysis)
				if o.OneAttrString(engine.DomainContext) != domaincontext {
					return false
				}
				return true
			}).Iterate(func(o *engine.Object) bool {

				grpsid := o.SID()
				if grpsid.IsNull() {
					return true
				}

				switch grpsid.RID() {
				case DOMAIN_USER_RID_ADMIN:
				case DOMAIN_USER_RID_KRBTGT:
				case DOMAIN_GROUP_RID_ADMINS:
				case DOMAIN_GROUP_RID_CONTROLLERS:
				case DOMAIN_GROUP_RID_SCHEMA_ADMINS:
				case DOMAIN_GROUP_RID_ENTERPRISE_ADMINS:
				case DOMAIN_GROUP_RID_READONLY_CONTROLLERS:
				case DOMAIN_ALIAS_RID_ADMINS:
				case DOMAIN_ALIAS_RID_ACCOUNT_OPS:
					if excluded_mask&1 != 0 {
						return true
					}
				case DOMAIN_ALIAS_RID_SYSTEM_OPS:
					if excluded_mask&2 != 0 {
						return true
					}
				case DOMAIN_ALIAS_RID_PRINT_OPS:
					if excluded_mask&4 != 0 {
						return true
					}
				case DOMAIN_ALIAS_RID_BACKUP_OPS:
					if excluded_mask&8 != 0 {
						return true
					}
				case DOMAIN_ALIAS_RID_REPLICATOR:
				default:
					// Not a protected group
					return true
				}

				// Only domain groups
				if grpsid.Component(2) != 21 && grpsid.Component(2) != 32 {
					ui.Debug().Msgf("RID match but not domain object for %v with SID %v", o.OneAttrString(engine.DistinguishedName), o.SID().String())
					return true
				}

				// Apply this edge
				adminsdholder.EdgeTo(o, activedirectory.EdgeOverwritesACL)
				o.EdgeIteratorRecursive(engine.In, engine.EdgeBitmap{}.Set(activedirectory.EdgeMemberOfGroup), true, func(source, target *engine.Object, edge engine.EdgeBitmap, depth int) bool {
					adminsdholder.EdgeTo(target, activedirectory.EdgeOverwritesACL)
					return true
				})
				return true
			})
			return true
		})
	},
		"AdminSDHolder rights propagation indicator",
		engine.BeforeMerge)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Add our known SIDs if they're missing
		for sid, name := range windowssecurity.KnownSIDs {
			binsid, err := windowssecurity.ParseStringSID(sid)
			if err != nil {
				ui.Fatal().Msgf("Problem parsing SID %v", sid)
			}
			if fo := FindWellKnown(ao, binsid); fo == nil {
				dn := "CN=" + name + ",CN=microsoft-builtin"
				ui.Debug().Msgf("Adding missing well known SID %v (%v) as %v", name, sid, dn)
				ao.Add(engine.NewObject(
					engine.DistinguishedName, engine.NewAttributeValueString(dn),
					engine.Name, engine.NewAttributeValueString(name),
					engine.ObjectSid, engine.NewAttributeValueSID(binsid),
					engine.ObjectClass, engine.NewAttributeValueString("person"), engine.NewAttributeValueString("user"), engine.NewAttributeValueString("top"),
					engine.Type, engine.NewAttributeValueString("Group"),
				))
			}
		}
	},
		"missing well-known SIDs",
		engine.BeforeMergeLow,
	)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Generate member of chains
		everyonesid, _ := windowssecurity.ParseStringSID("S-1-1-0")
		everyone := FindWellKnown(ao, everyonesid)
		if everyone == nil {
			ui.Fatal().Msgf("Could not locate Everyone, aborting - this should at least have been added during earlier preprocessing")
		}

		authenticateduserssid, _ := windowssecurity.ParseStringSID("S-1-5-11")
		authenticatedusers := FindWellKnown(ao, authenticateduserssid)
		if authenticatedusers == nil {
			ui.Fatal().Msgf("Could not locate Authenticated Users, aborting - this should at least have been added during earlier preprocessing")
		}

		authenticatedusers.EdgeTo(everyone, activedirectory.EdgeMemberOfGroup)

		ncname, netbiosname, dnsroot, domainsid, err := FindDomain(ao)
		if err != nil {
			ui.Fatal().Msgf("Could not get needed domain information (%v), aborting", err)
		}

		DCsyncObject, _ := ao.FindTwoOrAdd(
			engine.Type, engine.ObjectTypeCallableServicePoint.ValueString(),
			engine.Name, engine.NewAttributeValueString("DCsync"),
		)
		DCsyncObject.Tag("hvt")

		dnsroot = strings.ToLower(dnsroot)
		TrustMap.Store(TrustPair{
			SourceNCName:  ncname,
			SourceDNSRoot: dnsroot,
			SourceNetbios: netbiosname,
			SourceSID:     domainsid.String(),
		}, TrustInfo{})

		ao.Iterate(func(object *engine.Object) bool {
			if rid, ok := object.AttrInt(activedirectory.PrimaryGroupID); ok {
				sid := object.SID()
				if len(sid) > 8 {
					sidbytes := []byte(sid)
					binary.LittleEndian.PutUint32(sidbytes[len(sid)-4:], uint32(rid))
					primarygroup := ao.FindOrAddAdjacentSID(windowssecurity.SID(sidbytes), object)
					object.EdgeTo(primarygroup, activedirectory.EdgeMemberOfGroup)
				}
			}

			// Crude special handling for Everyone and Authenticated Users
			if object.SID().Components() == 7 && object.SID().StripRID() == domainsid && object.Type() != engine.ObjectTypeGroup {
				// if object.Type() == engine.ObjectTypeUser || object.Type() == engine.ObjectTypeComputer || object.Type() == engine.ObjectTypeManagedServiceAccount || object.Type() == engine.ObjectTypeGroupManagedServiceAccount {
				object.EdgeTo(authenticatedusers, activedirectory.EdgeMemberOfGroup)
			}

			if lastlogon, ok := object.AttrTime(activedirectory.LastLogonTimestamp); ok {
				object.Set(MetaLastLoginAge, engine.AttributeValueInt(int(time.Since(lastlogon)/time.Hour)))
			}
			if passwordlastset, ok := object.AttrTime(activedirectory.PwdLastSet); ok {
				object.Set(MetaPasswordAge, engine.AttributeValueInt(int(time.Since(passwordlastset)/time.Hour)))
			}
			if strings.Contains(strings.ToLower(object.OneAttrString(activedirectory.OperatingSystem)), "linux") {
				object.Tag("linux")
			}
			if strings.Contains(strings.ToLower(object.OneAttrString(activedirectory.OperatingSystem)), "windows") {
				object.Tag("windows")
			}
			if object.Attr(activedirectory.MSmcsAdmPwdExpirationTime).Len() > 0 {
				object.Tag("laps")
			}
			if uac, ok := object.AttrInt(activedirectory.UserAccountControl); ok {
				if uac&engine.UAC_TRUSTED_FOR_DELEGATION != 0 && uac&engine.UAC_NOT_DELEGATED == 0 {
					object.Tag("unconstrained")
				}
				if uac&engine.UAC_TRUSTED_TO_AUTH_FOR_DELEGATION != 0 {
					object.Tag("constrained")
				}
				if uac&engine.UAC_NOT_DELEGATED != 0 {
					object.Tag("nodelegation")
				}
				if uac&engine.UAC_WORKSTATION_TRUST_ACCOUNT != 0 {
					object.Tag("computer_account")
				}
				if uac&engine.UAC_SERVER_TRUST_ACCOUNT != 0 {
					object.Tag("domaincontroller_account")

					// All DCs are members of Enterprise Domain Controllers
					object.EdgeTo(ao.FindOrAddAdjacentSID(EnterpriseDomainControllers, object), activedirectory.EdgeMemberOfGroup)

					object.EdgeTo(DCsyncObject, activedirectory.EdgeCall)

					// Also they can DCsync because of this membership ... FIXME
				}

				var expired, disabled bool
				disabled = uac&engine.UAC_ACCOUNTDISABLE != 0
				if disabled {
					object.Tag("account_disabled")
				} else {
					object.Tag("account_enabled")
				}

				if uac&engine.UAC_LOCKOUT != 0 {
					object.Tag("account_locked")
				}

				if object.HasAttr(activedirectory.AccountExpires) {
					if exp, ok := object.Attr(activedirectory.AccountExpires).First().Raw().(time.Time); ok {
						if !exp.IsZero() && time.Now().After(exp) {
							object.Tag("account_expired")
							expired = true
						}
					}
				}

				if disabled || expired {
					object.Tag("account_inactive")
				} else {
					object.Tag("account_active")
				}
				if uac&engine.UAC_PASSWD_CANT_CHANGE != 0 {
					object.Tag("password_cant_change")
				}
				if uac&engine.UAC_DONT_EXPIRE_PASSWORD != 0 {
					object.Tag("password_never_expires")
				}
				if uac&engine.UAC_PASSWD_NOTREQD != 0 {
					object.Tag("password_not_required")
				}

				if uac&engine.UAC_SERVER_TRUST_ACCOUNT != 0 {
					// Domain Controller
					// find the machine object for this
					machine, found := ao.FindTwo(engine.Type, engine.NewAttributeValueString("Machine"),
						DomainJoinedSID, engine.NewAttributeValueSID(object.SID()))
					if !found {
						ui.Warn().Msgf("Can not find machine object for DC %v", object.DN())
					} else {
						machine.Tag("role_domaincontroller")
						machine.Tag("hvt")

						domainContext := object.OneAttr(engine.DomainContext)
						if domainContext == nil {
							ui.Fatal().Msgf("DomainController %v has no DomainContext attribute", object.DN())
						}

						if administrators, found := ao.FindTwo(engine.ObjectSid, engine.NewAttributeValueSID(windowssecurity.AdministratorsSID),
							engine.DomainContext, domainContext); found {
							administrators.EdgeTo(machine, activedirectory.EdgeLocalAdminRights)
						} else {
							ui.Warn().Msgf("Could not find Administrators group for %v", object.DN())
						}

						if remotedesktopusers, found := ao.FindTwo(engine.ObjectSid, engine.NewAttributeValueSID(windowssecurity.RemoteDesktopUsersSID),
							engine.DomainContext, domainContext); found {
							remotedesktopusers.EdgeTo(machine, activedirectory.EdgeLocalRDPRights)
						} else {
							ui.Warn().Msgf("Could not find Remote Desktop Users group for %v", object.DN())
						}

						if distributeddcomusers, found := ao.FindTwo(engine.ObjectSid, engine.NewAttributeValueSID(windowssecurity.DCOMUsersSID),
							engine.DomainContext, domainContext); found {
							distributeddcomusers.EdgeTo(machine, activedirectory.EdgeLocalDCOMRights)
						} else {
							ui.Warn().Msgf("Could not find DCOM Users group for %v", object.DN())
						}
					}
				}

				if object.HasAttrValue(activedirectory.PrimaryGroupID, engine.AttributeValueInt(521)) {
					// Read Only Domain Controller
					machine, found := ao.FindTwo(engine.Type, engine.NewAttributeValueString("Machine"),
						DomainJoinedSID, engine.NewAttributeValueSID(object.SID()))
					if !found {
						ui.Warn().Msgf("Can not find machine object for RODC %v", object.DN())
					} else {
						machine.Tag("role_readonly_domaincontroller")
						machine.Tag("hvt")
					}

					// Figure out what hashes this machine has cached - FIXME!

				}
			}

			if object.Type() == engine.ObjectTypeTrust {
				// http://www.frickelsoft.net/blog/?p=211
				var direction string
				dir, _ := object.AttrInt(activedirectory.TrustDirection)
				switch dir {
				case 0:
					direction = "disabled"
				case 1:
					direction = "incoming"
				case 2:
					direction = "outgoing"
				case 3:
					direction = "bidirectional"
				}

				attr, _ := object.AttrInt(activedirectory.TrustAttributes)

				partner := object.OneAttrString(activedirectory.TrustPartner)

				ui.Info().Msgf("Domain %v has a %v trust with %v", dnsroot, direction, partner)

				if dir&2 != 0 && attr&0x08 != 0 && attr&0x40 != 0 {
					ui.Info().Msgf("SID filtering is not enabled, so pwn %v and pwn this AD too", object.OneAttr(activedirectory.TrustPartner))
				}

				TrustMap.Store(TrustPair{
					SourceDNSRoot: dnsroot,
					TargetDNSRoot: partner,
				}, TrustInfo{
					Direction: TrustDirection(dir),
				})
			}

			/* else if object.HasAttrValue(engine.ObjectClass, "classSchema") {
				if u, ok := object.OneAttrRaw(engine.SchemaIDGUID).(uuid.UUID); ok {
					// ui.Debug().Msgf("Adding schema class %v %v", u, object.OneAttr(Name))
					engine.AllSchemaClasses[u] = object
				}
			}*/
			return true
		})
	},
		"Active Directory objects and metadata",
		engine.BeforeMergeHigh)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(object *engine.Object) bool {
			// We'll put the ObjectClass UUIDs in a synthetic attribute, so we can look it up later quickly (and without access to Objects)
			objectclasses := object.Attr(engine.ObjectClass)
			if objectclasses.Len() > 0 {
				guids := make([]engine.AttributeValue, 0, objectclasses.Len())
				objectclasses.Iterate(func(class engine.AttributeValue) bool {
					if oto, found := ao.Find(engine.LDAPDisplayName, class); found {
						if guid, ok := oto.OneAttr(activedirectory.SchemaIDGUID).(engine.AttributeValueGUID); !ok {
							ui.Debug().Msgf("%v", oto)
							ui.Fatal().Msgf("Could not translate SchemaIDGUID for class %v - I need a Schema to work properly", class)
						} else {
							guids = append(guids, guid)
						}
					} else {
						ui.Warn().Msgf("Could not resolve object class %v, perhaps you didn't get a dump of the schema?", class.String())
					}
					return true // continue
				})
				object.SetFlex(engine.ObjectClassGUIDs, guids)
			}

			// ObjectCategory handling
			var objectcategoryguid engine.AttributeValue
			var simple engine.AttributeValue

			objectcategoryguid = engine.NewAttributeValueGUID(engine.UnknownGUID)
			simple = engine.NewAttributeValueString("Unknown")

			typedn := object.OneAttr(engine.ObjectCategory)

			// Does it have one, and does it have a comma, then we're assuming it's not just something we invented
			if typedn != nil {
				if oto, found := ao.Find(engine.DistinguishedName, typedn); found {
					if _, ok := oto.OneAttrRaw(activedirectory.SchemaIDGUID).(uuid.UUID); ok {
						objectcategoryguid = oto.OneAttr(activedirectory.SchemaIDGUID)
						simple = oto.OneAttr(activedirectory.Name)
					} else {
						ui.Error().Msgf("Could not translate SchemaIDGUID for %v", typedn)
					}
				} else {
					ui.Error().Msgf("Could not resolve object category %v, perhaps you didn't get a dump of the schema?", typedn)
				}
			}

			object.SetFlex(
				engine.ObjectCategoryGUID, objectcategoryguid,
				engine.Type, simple,
			)
			return true
		})
	},
		"Set type (for Type call) to Active Directory objects",
		engine.BeforeMergeLow,
	)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(object *engine.Object) bool {
			if object.SID().Component(2) == 21 && object.SID().RID() == 525 { // "Protected Users"
				object.EdgeIteratorRecursive(engine.In, engine.EdgeBitmap{}.Set(activedirectory.EdgeMemberOfGroup), true, func(source, member *engine.Object, edge engine.EdgeBitmap, depth int) bool {
					if member.Type() == engine.ObjectTypeComputer || member.Type() == engine.ObjectTypeUser {
						member.Tag("protected_user")
					}
					return true
				})
			}
			return true
		})
	},
		"Protected users meta attribute",
		engine.BeforeMerge,
	)

	// Loader.AddProcessor(func(ao *engine.Objects) {
	// 	// Find all the DomainDNS objects, and find the domain object
	// 	domains := make(map[string]windowssecurity.SID)

	// 	domaindnsobjects, found := ao.FindMulti(engine.ObjectClass, engine.NewAttributeValueString("domainDNS"))

	// 	if !found {
	// 		ui.Error().Msg("Could not find any domainDNS objects")
	// 	}

	// 	domaindnsobjects.Iterate(func(domaindnsobject *engine.Object) bool {
	// 		domainSID, sidok := domaindnsobject.OneAttrRaw(activedirectory.ObjectSid).(windowssecurity.SID)
	// 		dn := domaindnsobject.OneAttrString(activedirectory.DistinguishedName)
	// 		if sidok {
	// 			domains[dn] = domainSID
	// 		}
	// 		return true
	// 	})

	// 	ao.Iterate(func(o *engine.Object) bool {
	// 		if o.HasAttr(engine.ObjectSid) && o.SID().Component(2) == 21 && !o.HasAttr(engine.DistinguishedName) && o.HasAttr(engine.DomainContext) {
	// 			// An unknown SID, is it ours or from another domain?
	// 			ourDomainDN := o.OneAttrString(engine.DomainContext)
	// 			ourDomainSid, domainfound := domains[ourDomainDN]
	// 			if !domainfound {
	// 				return true
	// 			}

	// 			if o.SID().StripRID() == ourDomainSid {
	// 				// ui.Debug().Msgf("Found a 'dangling' local SID object %v. This is either a SID from a deleted object (most likely) or hardened objects that are not readable with the account used to dump data.", o.SID())
	// 			} else {
	// 				// ui.Debug().Msgf("Found a 'lost' foreign SID object %v, adding it as a synthetic Foreign-Security-Principal", o.SID())
	// 				o.SetFlex(
	// 					engine.DistinguishedName, engine.NewAttributeValueString(o.SID().String()+",CN=ForeignSecurityPrincipals,"+ourDomainDN),
	// 					engine.ObjectCategorySimple, "Foreign-Security-Principal",
	// 					engine.DataLoader, "Autogenerated",
	// 				)
	// 			}
	// 		}
	// 		return true
	// 	})
	// },
	// 	"Creation of synthetic Foreign-Security-Principal objects",
	// 	engine.AfterMergeLow)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(machine *engine.Object) bool {
			// Only for machines, you can't really pwn users this way
			if machine.Type() != ObjectTypeMachine {
				return true
			}

			// Find the computer AD object if any
			var computer *engine.Object
			machine.Edges(engine.Out).Range(func(target *engine.Object, edge engine.EdgeBitmap) bool {
				if edge.IsSet(EdgeAuthenticatesAs) && target.Type() == engine.ObjectTypeComputer {
					computer = target
					return false //break
				}
				return true
			})

			if computer == nil {
				ui.Warn().Msgf("Machine without computer account: %v", machine.Label())
				return true // continue
			}

			// Find all perent containers with GP links
			var hasparent bool

			// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpol/5c7ecdad-469f-4b30-94b3-450b7fff868f
			allowEnforcedGPOsOnly := false

			currentObject := computer
			var iteration int
			for {
				iteration++
				potentialParent := currentObject.Parent()
				if potentialParent != nil && potentialParent.DN() != "" && strings.HasSuffix(currentObject.DN(), potentialParent.DN()) {
					// It's usable
					currentObject = potentialParent
				} else {
					// Fall back to old slow method of looking at DNs
					currentObject, hasparent = ao.DistinguishedParent(currentObject)
					if !hasparent {
						break
					}
				}

				var gpcachelinks engine.AttributeValues
				var found bool
				if gpcachelinks, found = currentObject.Get(GPLinkCache); !found {
					// the hard way
					var gpcachelinks engine.AttributeValues

					gplinks := strings.Trim(currentObject.OneAttrString(activedirectory.GPLink), " ")
					if len(gplinks) != 0 {
						// ui.Debug().Msgf("GPlink for %v on container %v: %v", o.DN(), p.DN(), gplinks)
						if !strings.HasPrefix(gplinks, "[") || !strings.HasSuffix(gplinks, "]") {
							ui.Error().Msgf("Error parsing gplink on %v: %v", computer.DN(), gplinks)
						} else {
							links := strings.Split(gplinks[1:len(gplinks)-1], "][")

							var collecteddata engine.AttributeValues
							for _, link := range links {
								linkinfo := strings.Split(link, ";")
								if len(linkinfo) != 2 {
									ui.Error().Msgf("Error parsing gplink on %v: %v", computer.DN(), gplinks)
									continue
								}
								linkedgpodn := linkinfo[0][7:] // strip LDAP:// prefix and link to this

								gpo, found := ao.Find(engine.DistinguishedName, engine.NewAttributeValueString(linkedgpodn))
								if !found {
									if _, warned := warnedgpos[linkedgpodn]; !warned {
										warnedgpos[linkedgpodn] = struct{}{}
										ui.Warn().Msgf("Object linked to GPO that is not found %v: %v", computer.DN(), linkedgpodn)
									}
								} else {
									linktype, _ := strconv.ParseInt(linkinfo[1], 10, 64)
									collecteddata = append(collecteddata, engine.AttributeValueObject{
										Object: gpo,
									}, engine.AttributeValueInt(linktype))
								}
							}
							gpcachelinks = collecteddata
						}
					}
					currentObject.Set(GPLinkCache, gpcachelinks...)
				}

				// cached or generated - pairwise pointer to gpo object and int
				for i := 0; i < gpcachelinks.Len(); i += 2 {
					gpo := gpcachelinks[i].Raw().(*engine.Object)
					gpLinkOptions := gpcachelinks[i+1].Raw().(int64)
					// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-gpol/08090b22-bc16-49f4-8e10-f27a8fb16d18
					if gpLinkOptions&0x01 != 0 {
						// GPO link is disabled
						continue
					}
					if allowEnforcedGPOsOnly && gpLinkOptions&0x02 == 0 {
						// Enforcement required, but this is not an enforced GPO
						continue
					}
					gpo.EdgeTo(machine, activedirectory.EdgeAffectedByGPO)
				}
				gpoptions := currentObject.OneAttrString(activedirectory.GPOptions)
				if gpoptions == "1" {
					// inheritance is blocked, so let's not forget that when moving up
					allowEnforcedGPOsOnly = true
				}
			}
			return true
		})
	},
		"Machines affected by a GPO",
		engine.AfterMergeLow,
	)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			if o.HasAttr(engine.ObjectSid) && !o.HasAttr(engine.DisplayName) {
				if name, found := windowssecurity.KnownSIDs[o.SID().String()]; found {
					o.SetFlex(engine.DisplayName, name)
				}
			}
			return true
		})
	},
		"Adding displayName to Well-Known SID objects that are missing them",
		engine.AfterMergeLow)

	// CREATOR_OWNER is a template for new objects, so this was totally wrong
	/*
		Loader.AddProcessor(func(ao *engine.Objects) {
			creatorowner, found := ao.Find(engine.ObjectSid, engine.AttributeValueSID(windowssecurity.CreatorOwnerSID))
			if !found {
				ui.Warn().Msg("Could not find Creator Owner Well Known SID. Not doing post-merge fixup")
				return
			}

			for target, edges := range creatorowner.CanPwn {
				// ACL grants CreatorOwnerSID something - so let's find the owner and give them the permissions
				if sd, err := target.SecurityDescriptor(); err == nil {
					if sd.Owner != windowssecurity.BlankSID {
						if realowners, found := ao.FindMulti(engine.ObjectSid, engine.AttributeValueSID(sd.Owner)); found {
							for _, realo := range realowners {
								if realo.Type() == engine.ObjectTypeForeignSecurityPrincipal || realo.Type() == engine.ObjectTypeOther {
									// Skip this
									continue
								}

								// Link real target
								realo.CanPwn[target] = realo.CanPwn[target].Merge(edges)
								target.PwnableBy[realo] = target.PwnableBy[realo].Merge(edges)

								// Unlink creatorowner
								delete(creatorowner.CanPwn, target)
								delete(target.PwnableBy, creatorowner)
							}
						}
					}
				}
			}
		},
			"CreatorOwnerSID resolution fixup",
			engine.BeforeMerge,
		)
	*/

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(object *engine.Object) bool {
			// Object that is member of something
			object.Attr(activedirectory.MemberOf).Iterate(func(memberof engine.AttributeValue) bool {
				group, found := ao.Find(engine.DistinguishedName, memberof)
				if !found {
					var sid engine.AttributeValue
					if stringsid, _, found := strings.Cut(memberof.String(), ",CN=ForeignSecurityPrincipals,"); found {
						// We can figure out what the SID is
						if c, err := windowssecurity.ParseStringSID(stringsid); err == nil {
							sid = engine.NewAttributeValueSID(c)
						}
						ui.Info().Msgf("Missing Foreign-Security-Principal: %v is a member of %v, which is not found - adding enhanced synthetic group", object.DN(), memberof)
					} else {
						ui.Warn().Msgf("Possible hardening? %v is a member of %v, which is not found - adding synthetic group. Your analysis will be degraded, try dumping with Domain Admin rights.", object.DN(), memberof)
					}
					group = engine.NewObject(
						engine.IgnoreBlanks,
						engine.DistinguishedName, memberof,
						engine.Type, engine.NewAttributeValueString("Group"),
						engine.ObjectClass, engine.NewAttributeValueString("top"), engine.NewAttributeValueString("group"),
						engine.Name, engine.NewAttributeValueString("Synthetic group "+memberof.String()),
						engine.Description, engine.NewAttributeValueString("Synthetic group"),
						engine.ObjectSid, sid,
						engine.DataLoader, engine.NewAttributeValueString("Autogenerated"),
					)
					ao.Add(group)
				}
				object.EdgeTo(group, activedirectory.EdgeMemberOfGroup)
				return true
			})

			// Group that contains members
			object.Attr(activedirectory.Member).Iterate(func(member engine.AttributeValue) bool {
				memberobject, found := ao.Find(engine.DistinguishedName, member)
				if !found {
					if stringsid, _, found := strings.Cut(member.String(), ",CN=ForeignSecurityPrincipals,"); found {
						// We can figure out what the SID is
						stringsid, _, _ = strings.Cut(stringsid[3:], "\\") // remote CN= and \=ACNF:guid

						if sid, err := windowssecurity.ParseStringSID(stringsid); err == nil {
							memberobject = ao.FindOrAddAdjacentSID(sid, object)
						} else {
							ui.Warn().Msgf("Could not extract SID from Foreign-Security-Principal %v: %v", member.String(), err)
						}
					}
					if memberobject == nil {
						ui.Warn().Msgf("Possible hardening? %v is a member of %v, which is not found - adding synthetic member. Your analysis will be degraded, try dumping with Domain Admin rights.", member, object.DN())
						memberobject, _ = ao.FindOrAdd(engine.DistinguishedName, member,
							engine.DataLoader, "Autogenerated",
						)
					}
				}
				memberobject.EdgeTo(object, activedirectory.EdgeMemberOfGroup)
				return true
			})
			return true
		})
	},
		"MemberOf and Member resolution",
		engine.AfterMergeLow,
	)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		ao.Iterate(func(o *engine.Object) bool {
			// Only for containers and org units
			if o.Type() != engine.ObjectTypeUser {
				return true
			}

			sd, err := o.SecurityDescriptor()
			if err != nil {
				return true
			}
			for index, acl := range sd.DACL.Entries {
				if sd.DACL.IsObjectClassAccessAllowed(index, o, engine.RIGHT_DS_WRITE_PROPERTY, AttributeUserAccountControlGUID, ao) {
					ao.FindOrAddAdjacentSID(acl.SID, o).EdgeTo(o, activedirectory.EdgeWriteUserAccountControl)
				}
			}
			return true
		})
	}, "Permissions that lets someone modify userAccountControl", engine.BeforeMergeFinal)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		edgematch := engine.EdgeBitmap{}.Set(activedirectory.EdgeMemberOfGroup).Set(activedirectory.EdgeForeignIdentity)
		cacheMap := xsync.NewTypedMapOf[*engine.Object, []engine.AttributeValue](func(o *engine.Object) uint64 {
			return uint64(uintptr(unsafe.Pointer(o)))
		})

		ao.IterateParallel(func(o *engine.Object) bool {
			// Search from all groups towards incoming memberships
			o.EdgeIteratorRecursive(engine.Out, edgematch, true, func(member, memberof *engine.Object, edge engine.EdgeBitmap, depth int) bool {
				if memberOfIndirects, found := cacheMap.Load(memberof); found {
					o.Add(MemberOfIndirect, memberOfIndirects...)
					return false // don't go deeper
				}
				if depth > 1 && memberof.Type() == engine.ObjectTypeGroup {
					member.EdgeTo(memberof, activedirectory.EdgeMemberOfGroupIndirect)

					dn := memberof.Attr(engine.DistinguishedName)
					if dn.First() != nil {
						o.Add(MemberOfIndirect, dn.First())
					}
				}
				return true // deeper!!
			})

			// Populate cache
			memberOfIndirect, _ := o.Get(MemberOfIndirect)
			cacheMap.Store(o, memberOfIndirect)

			return true
		}, 0)
	},
		"MemberOfIndirect resolution",
		engine.AfterMerge,
	)

	LoaderID.AddProcessor(
		func(ao *engine.Objects) {
			ao.Iterate(func(enrollementService *engine.Object) bool {
				if enrollementService.Type() == engine.ObjectTypePKIEnrollmentService {
					var ca *engine.Object
					var found bool
					if cadns := enrollementService.OneAttr(activedirectory.DNSHostName); cadns != nil {
						// find the CA machine object
						if ca, found = ao.FindTwo(
							engine.Type, ObjectTypeMachine.ValueString(),
							activedirectory.DNSHostName, cadns,
						); found {
							ca.Tag("role_certificate_authority")
							ca.Tag("hvt")
						} else {
							ui.Warn().Msgf("Couldn't locate dnsHostName %v acting as enrollmentservice", cadns)
						}
					}

					// Templates that is offered for enrollment
					enrollementService.Attr(CertificateTemplates).Iterate(func(templatename engine.AttributeValue) bool {

						templates, found := ao.FindTwoMulti(engine.Name, templatename,
							engine.ObjectClass, engine.NewAttributeValueString("pKICertificateTemplate"))

						if found {
							alreadyset := false
							templates.Iterate(func(template *engine.Object) bool {
								if !engine.CompareAttributeValues(template.OneAttr(engine.DomainContext), enrollementService.OneAttr(engine.DomainContext)) {
									return true // continue
								}

								if alreadyset {
									ui.Warn().Msgf("Found multiple templates for %s", templatename)
								}

								template.SetFlex(
									PublishedBy, engine.NewAttributeValueString(enrollementService.DN()),
									PublishedByDnsHostName, enrollementService.Attr(activedirectory.DNSHostName),
								)

								template.Tag("published")

								// classify the template as ESC1 - 11

								alreadyset = true
								return true
							})
							if !alreadyset {
								ui.Warn().Msgf("Found no matching template for %s", templatename)
							}
						} else {
							ui.Warn().Msgf("Template %s not found", templatename)
						}
						return true
					})
				}
				return true
			})
		},
		"Certificate template publishing status",
		engine.AfterMerge,
	)

	/*
		LoaderID.AddProcessor(func(ao *engine.Objects) {
			ao.Iterate(func(enrollmentservice *engine.Object) bool {
				if enrollementService.Type() == engine.ObjectTypePKIEnrollmentService {
				var ca *engine.Object
				var found bool
				if cadns := enrollementService.OneAttr(activedirectory.DNSHostName); cadns != nil {
					// find the CA machine object
					if ca, found = ao.FindTwo(
						engine.Type, ObjectTypeMachine.ValueString(),
						activedirectory.DNSHostName, cadns,
					); !found {
						return true
					}
				}

				// we have a ca

			})
		})
	*/

	/*
		Loader.AddProcessor(func(ao *engine.Objects) {
			ao.Filter(func(o *engine.Object) bool {
				return o.Type() == engine.ObjectTypeForeignSecurityPrincipal
			}).Iterate(func(foreign *engine.Object) bool {
				sid := foreign.SID()
				if sid.IsNull() {
					ui.Error().Msgf("Found a foreign security principal with no SID %v", foreign.Label())
					return true
				}
				if sid.Component(2) == 21 {
					if sources, found := ao.FindMulti(engine.ObjectSid, engine.AttributeValueSID(sid)); found {
						sources.Iterate(func(source *engine.Object) bool {
							if source.Type() != engine.ObjectTypeForeignSecurityPrincipal {
								source.EdgeToEx(foreign, activedirectory.EdgeForeignIdentity, true)
							}
							return true
						})
					}
				} else {
					ui.Warn().Msgf("Found a foreign security principal %v with an non type 21 SID %v", foreign.DN(), sid.String())
				}
				return true
			})
		}, "Link foreign security principals to their native objects",
			engine.AfterMerge,
		)
	*/

	type sidinfo struct {
		domainContext string
	}

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		// Find all domains, save info so we can see if an object is "local" or not
		sidmap := make(map[windowssecurity.SID]sidinfo)
		ao.Filter(func(o *engine.Object) bool {
			return o.HasAttr(activedirectory.ObjectSid) && o.Type() == engine.ObjectTypeDomainDNS
		}).Iterate(func(domain *engine.Object) bool {
			sid := domain.SID()
			domainContext := domain.OneAttrString(engine.DomainContext)
			sidmap[sid] = sidinfo{
				domainContext: domainContext,
			}
			return true
		})

		ao.Filter(func(o *engine.Object) bool {
			return o.HasAttr(activedirectory.ObjectSid)
		}).Iterate(func(object *engine.Object) bool {
			sid := object.SID()
			if object.HasAttr(engine.DomainContext) {
				domainContext := object.OneAttrString(engine.DomainContext)
				domaininfo, found := sidmap[sid.StripRID()]
				if found && domaininfo.domainContext != domainContext {
					// it's foreign, find the local one
					nativeObjects, found := ao.FindTwoMulti(
						engine.ObjectSid, engine.NewAttributeValueSID(sid),
						engine.DomainContext, engine.NewAttributeValueString(domainContext),
					)
					if found {
						nativeobject := nativeObjects.First()
						nativeobject.EdgeTo(object, activedirectory.EdgeForeignIdentity)
						// Inherit the type from the original
						if !object.HasAttr(activedirectory.Type) {
							object.SetFlex(activedirectory.Type, nativeobject.Attr(activedirectory.Type))
						}
					}
				}
			}
			return true
		})
	}, "Link foreign security principals to their native objects",
		engine.AfterMergeLow,
	)

	LoaderID.AddProcessor(func(ao *engine.Objects) {
		var warnlines int
		ao.Filter(func(o *engine.Object) bool {
			return o.Type() == engine.ObjectTypeGroupPolicyContainer
		}).Iterate(func(gpo *engine.Object) bool {
			gpo.Edges(engine.In).Range(func(group *engine.Object, methods engine.EdgeBitmap) bool {
				groupname := group.OneAttrString(engine.SAMAccountName)
				if strings.Contains(groupname, "%") {
					// Lowercase for ease
					groupname := strings.ToLower(groupname)

					// It has some sort of % variable in it, let's go
					gpo.Edges(engine.Out).Range(func(affected *engine.Object, amethods engine.EdgeBitmap) bool {
						if amethods.IsSet(activedirectory.EdgeAffectedByGPO) && affected.Type() == engine.ObjectTypeComputer {
							netbiosdomain, computername, found := strings.Cut(affected.OneAttrString(engine.DownLevelLogonName), "\\")
							if !found {
								ui.Error().Msgf("Could not parse downlevel logon name %v", affected.OneAttrString(engine.DownLevelLogonName))
								return true //continue
							}
							computername = strings.TrimRight(computername, "$")

							realgroup := groupname
							realgroup = strings.Replace(realgroup, "%computername%", computername, -1)
							realgroup = strings.Replace(realgroup, "%domainname%", netbiosdomain, -1)
							realgroup = strings.Replace(realgroup, "%domain%", netbiosdomain, -1)

							var targetgroups engine.ObjectSlice

							if !strings.Contains(realgroup, "\\") {
								realgroup = netbiosdomain + "\\" + realgroup
							}
							targetgroups, _ = ao.FindMulti(
								engine.DownLevelLogonName, engine.NewAttributeValueString(realgroup),
							)

							if targetgroups.Len() == 0 {
								if warnlines < 10 {
									ui.Warn().Msgf("Could not find group %v", realgroup)
								}
								warnlines++
							} else if targetgroups.Len() == 1 {
								for _, edge := range methods.Edges() {
									targetgroups.First().EdgeToEx(affected, edge, true)
								}
							} else {
								ui.Warn().Msgf("Found multiple groups for %v: %v", realgroup, targetgroups)
								targetgroups.Iterate(func(targetgroup *engine.Object) bool {
									ui.Warn().Msgf("Target: %v", targetgroup.DN())
									return true
								})
							}
						}
						return true
					})
				}
				return true
			})
			return true
		})
		if warnlines > 0 {
			ui.Warn().Msgf("%v groups could not be resolved, this could affect analysis results", warnlines)
		}

	}, "Resolve expanding environment variables in group names to real names from GPOs",
		engine.AfterMerge,
	)
}
