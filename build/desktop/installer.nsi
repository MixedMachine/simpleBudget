!define APP_NAME "simple-budget-app"
!define APP_VERSION "0.3.0"
!define APP_OUT_DIR "SimpleBudgetApp"

; Set the name of the installer and the output directory
Outfile "${APP_NAME}-${APP_VERSION}-installer.exe"

; Set the default installation directory
InstallDir $PROGRAMFILES\${APP_OUT_DIR}

ShowInstDetails show

UninstPage uninstConfirm
LicenseData "..\..\LICENSE.txt"

; Start the installation process
Section "install"
    SetOutPath $INSTDIR
    File "..\..\bin\desktop\${APP_NAME}.exe"
    File "..\..\assets\icon.png"

    ; Create shortcuts in the Start Menu
    CreateDirectory "$SMPROGRAMS\${APP_OUT_DIR}"
    CreateShortCut "$SMPROGRAMS\${APP_OUT_DIR}\Simple Budget App.lnk" "$INSTDIR\${APP_NAME}.exe"

    WriteUninstaller "$INSTDIR\uninstall.exe"

SectionEnd

; Create the uninstaller
Section "uninstall"
    ; Remove the installed files
    Delete $INSTDIR\YourBudgetApp.exe
    Delete $INSTDIR\icon.png

    ; Remove shortcuts from the Start Menu
    Delete "$SMPROGRAMS\${APP_OUT_DIR}\Simple Budget App.lnk"

    RMDir "$SMPROGRAMS\${APP_OUT_DIR}"

    ; Prompt the user to delete the installation directory
    RMDir $INSTDIR

SectionEnd
