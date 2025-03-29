package engine

const (
	UAC_SCRIPT                         = 0x0001
	UAC_ACCOUNTDISABLE                 = 0x0002
	UAC_HOMEDIR_REQUIRED               = 0x0008
	UAC_LOCKOUT                        = 0x0010
	UAC_PASSWD_NOTREQD                 = 0x0020
	UAC_PASSWD_CANT_CHANGE             = 0x0040
	UAC_ENCRYPTED_TEXT_PWD_ALLOWED     = 0x0080
	UAC_TEMP_DUPLICATE_ACCOUNT         = 0x0100
	UAC_NORMAL_ACCOUNT                 = 0x0200
	UAC_INTERDOMAIN_TRUST_ACCOUNT      = 0x0800
	UAC_WORKSTATION_TRUST_ACCOUNT      = 0x1000
	UAC_SERVER_TRUST_ACCOUNT           = 0x2000
	UAC_DONT_EXPIRE_PASSWORD           = 0x10000
	UAC_MNS_LOGON_ACCOUNT              = 0x20000
	UAC_SMARTCARD_REQUIRED             = 0x40000
	UAC_TRUSTED_FOR_DELEGATION         = 0x80000
	UAC_NOT_DELEGATED                  = 0x100000
	UAC_USE_DES_KEY_ONLY               = 0x200000
	UAC_DONT_REQ_PREAUTH               = 0x400000
	UAC_PASSWORD_EXPIRED               = 0x800000
	UAC_TRUSTED_TO_AUTH_FOR_DELEGATION = 0x1000000
	UAC_PARTIAL_SECRETS_ACCOUNT        = 0x04000000

	RIGHT_GENERIC_READ Mask = RIGHT_READ_CONTROL | RIGHT_DS_LIST_CONTENTS | RIGHT_DS_READ_PROPERTY | RIGHT_DS_LIST_OBJECT /*
		** Mask value is not stored in AD but deduced from mask bits combined **
		RIGHT_GENERIC_READ = 0x80000000 /*
			The right to read permissions and all properties of the object, and list the contents of the
			object in the case of containers.

			Equivalent to:RIGHT_READ_CONTROL | RIGHT_DS_LIST_CONTENTS | RIGHT_DS_READ_PROPERTY | RIGHT_DS_LIST_OBJECT */

	RIGHT_GENERIC_WRITE = RIGHT_READ_CONTROL | RIGHT_DS_WRITE_PROPERTY | RIGHT_DS_WRITE_PROPERTY_EXTENDED /*
		** Mask value is not stored in AD but deduced from mask bits combined **
		RIGHT_GENERIC_WRITE = 0x40000000 /*
			Includes the right to read permissions on the object, and the right to write all the properties
			on the object.

			Equivalent to: RIGHT_READ_CONTROL | RIGHT_DS_WRITE_PROPERTY | RIGHT_DS_WRITE_PROPERTY_EXTENDED */

	RIGHT_GENERIC_EXECUTE = RIGHT_READ_CONTROL | RIGHT_DS_LIST_CONTENTS /*
		** Mask value is not stored in AD but deduced from mask bits combined **
		RIGHT_GENERIC_EXECUTE = 0x20000000 /*
			The right to read permissions/list the contents of a container object.

			Equivalent to: RIGHT_READ_CONTROL | RIGHT_DS_LIST_CONTENTS */
	RIGHT_GENERIC_ALL = RIGHT_DELETE | RIGHT_READ_CONTROL | RIGHT_WRITE_DACL | RIGHT_WRITE_OWNER | RIGHT_DS_CREATE_CHILD | RIGHT_DS_DELETE_CHILD | RIGHT_DS_DELETE_TREE | RIGHT_DS_READ_PROPERTY | RIGHT_DS_WRITE_PROPERTY | RIGHT_DS_LIST_CONTENTS | RIGHT_DS_LIST_OBJECT | RIGHT_DS_CONTROL_ACCESS | RIGHT_DS_WRITE_PROPERTY_EXTENDED /*
		** Mask value is not stored in AD but deduced from mask bits combined **
		RIGHT_GENERIC_ALL = 0x10000000 /*
			The right to create/delete child objects, read/write all properties, see any child objects, add and remove the object,
			and read/write with an extended right.

			Equivalent to: RIGHT_DELETE |  RIGHT_READ_CONTROL | RIGHT_WRITE_DACL | RIGHT_WRITE_OWNER | RIGHT_DS_CREATE_CHILD | RIGHT_DS_DELETE_CHILD | RIGHT_DS_DELETE_TREE | RIGHT_DS_READ_PROPERTY | RIGHT_DS_WRITE_PROPERTY | RIGHT_DS_LIST_CONTENTS | RIGHT_DS_LIST_OBJECT | RIGHT_DS_CONTROL_ACCESS | RIGHT_DS_WRITE_PROPERTY_EXTENDED)
	*/

	RIGHT_SYNCRONIZE  = 0x00100000
	RIGHT_WRITE_OWNER = 0x00080000 /*
		The right to modify the owner section of the security descriptor. Of note, a user with this right can only change the owner to themselves
		-ownership cannot be transferred to other userswith only this right.*/
	RIGHT_WRITE_DACL = 0x00040000 /*
		The right to modify the DACL for the object. */
	RIGHT_READ_CONTROL = 0x00020000 /*
		The right to read alldata from the security descriptor except the SACL. */
	RIGHT_DELETE = 0x00010000 /*
		The right to delete the object. */

	RIGHT_DS_VOODOO_BIT = 0x00001000 /* No clue - see https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-crtd/4be42fa6-c421-4763-890b-07a9ab5a319d for second option */

	RIGHT_DS_CONTROL_ACCESS = 0x00000100 /*
		A specific control access right (if the ObjectType GUID refers to an extended right registered in the forest schema)
		or the right to read a confidential property (if the ObjectType GUID refers to a confidential property).
		If the GUID is not present, then all extended rights are granted */
	RIGHT_DS_LIST_OBJECT = 0x00000080 /*
		The right to list an object. If the user does not have this right and also doesn’t have the
		RIGHT_DS_LIST_CONTENTS right on the object's parent container then the object is hidden from the user. */
	RIGHT_DS_DELETE_TREE = 0x00000040 /*
		The right to perform a delete-tree operation. */
	RIGHT_DS_WRITE_PROPERTY = 0x00000020 /*
		The right to write one or more properties of the object specified by the ObjectType GUID.
		If the ObjectType GUID is not present or is all 0s, then the right to write all properties is granted. */
	RIGHT_DS_READ_PROPERTY = 0x00000010 /*
		The right to read one or more properties of the object specified by the ObjectType GUID.
		If the ObjectType GUID is not present or is all 0s, then the right to read all properties is granted.	*/
	RIGHT_DS_WRITE_PROPERTY_EXTENDED = 0x00000008 /*
		The right to execute a validated write access right. AKA DsSelf */
	RIGHT_DS_LIST_CONTENTS = 0x00000004 /*
		The right to list all child objects of the object, if the object is a type of container. */
	RIGHT_DS_DELETE_CHILD = 0x00000002 /*
		The right to delete child objects of the object, if the object is a type of container.
		If the ObjectType contains a GUID, the GUID will reference the type of child object that can be deleted. */
	RIGHT_DS_CREATE_CHILD = 0x00000001 /*
		The right to create child objects under the object, if the object is a type of container.
		If the ObjectType contains a GUID, the GUID will reference the type of child object that can be created. */

)
