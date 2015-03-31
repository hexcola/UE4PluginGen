package main

import(
	"fmt"
	"bufio"
	"os"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
	)


type Module struct{
	Name string
	pluginType string `json:"Type"`
}


type Plugin struct{
	FileVersion int

	// friendly name
	FriendlyName string
	Version int
	VersionName string
	CreatedBy string
	CreatedByURL string
	EngineVersion string
	Description string
	Category string
	EnabledByDefault bool

	Modules []Module
}

type Config struct{
	ProjectPath string
	ThePlugin Plugin
}


func main(){

	dat, err := ioutil.ReadFile("config.ini")
	check(err)
	
	var config Config
	myJson := string(dat)

	json.Unmarshal([]byte(myJson), &config)


	plugin := config.ThePlugin


	
	// convert plugin to json

	pluginJson, err := json.Marshal(plugin)
	check(err)

	fmt.Println(string(pluginJson))




	pluginsFolder := config.ProjectPath + string(filepath.Separator) + "Plugins"
	pluginFolder := pluginsFolder+ string(filepath.Separator) + plugin.FriendlyName
	resourcesFolder := pluginFolder + string(filepath.Separator) + "Resources"
	
	sourceFolder := pluginFolder + string(filepath.Separator) + "Source"
	sourcePluginFolder := sourceFolder + string(filepath.Separator) + plugin.FriendlyName
	classesFolder := sourcePluginFolder + string(filepath.Separator) + "Classes"
	privateFolder := sourcePluginFolder + string(filepath.Separator) + "Private"
	publicFolder := sourcePluginFolder + string(filepath.Separator) + "Public"

	structure := []string {pluginsFolder, pluginFolder, resourcesFolder, sourceFolder, sourcePluginFolder, classesFolder, privateFolder, publicFolder}

	for _,v := range structure{
		fmt.Println(v)
		os.Mkdir(v, 0777)
	}


	// generate *.uplugin file
	upluginPath := pluginFolder + string(filepath.Separator) + plugin.FriendlyName + ".uplugin"
	err = ioutil.WriteFile(upluginPath, pluginJson, 0644)
	check(err)

	

	// generate *.Build.cs file
	buildCsFilePath := sourcePluginFolder + string(filepath.Separator) + plugin.FriendlyName + ".Build.cs"
	genBuildCsFile(buildCsFilePath, plugin)

	// generate *.I<PluginName>.h
	iPluginNamePath := publicFolder + string(filepath.Separator) + "I" + plugin.FriendlyName + ".h"
	genIPluginNameFile(iPluginNamePath, plugin.FriendlyName)

	// generate *.

}

func genBuildCsFile(path string, plugin Plugin){
	buildCsFile := []byte (`// Copyright 1998-2015 Epic Games, Inc. All Rights Reserved.

namespace UnrealBuildTool.Rules
{
	public class ` + plugin.FriendlyName +` : ModuleRules
	{
		public `+ plugin.FriendlyName +` (TargetInfo Target)
		{
			PublicIncludePaths.AddRange(
				new string[] {
					// ... add public include paths required here ...
				}
				);

			PrivateIncludePaths.AddRange(
				new string[] {
					//"Runtime/UEduino/Private",
					// ... add other private include paths required here ...
				}
				);

			PublicDependencyModuleNames.AddRange(
				new string[]
				{
					"Core",
					// ... add other public dependencies that you statically link with here ...
				}
				);

			PrivateDependencyModuleNames.AddRange(
				new string[]
				{
					// ... add private dependencies that you statically link with here ...
				}
				);

			DynamicallyLoadedModuleNames.AddRange(
				new string[]
				{
					// ... add any modules that your module loads dynamically here ...
				}
				);
		}
	}
}`)


	err := ioutil.WriteFile(path, buildCsFile, 0644)
	check(err)
}


func genIPluginNameFile(path string, pluginName string){
	IPlugingName := []byte(`// Copyright 1998-2015 Epic Games, Inc. All Rights Reserved.

#pragma once

#include "ModuleManager.h"


/**
 * The public interface to this module.  In most cases, this interface is only public to sibling modules 
 * within this plugin.
 */
class I`+ pluginName+` : public IModuleInterface
{

public:

	/**
	 * Singleton-like access to this module's interface.  This is just for convenience!
	 * Beware of calling this during the shutdown phase, though.  Your module might have been unloaded already.
	 *
	 * @return Returns singleton instance, loading the module on demand if needed
	 */
	static inline I`+pluginName+`& Get()
	{
		return FModuleManager::LoadModuleChecked< I`+pluginName+` >( "`+pluginName+`" );
	}

	/**
	 * Checks to see if this module is loaded and ready.  It is only valid to call Get() if IsAvailable() returns true.
	 *
	 * @return True if the module is loaded and ready to use
	 */
	static inline bool IsAvailable()
	{
		return FModuleManager::Get().IsModuleLoaded( "`+pluginName+`" );
	}
};`)

	err := ioutil.WriteFile(path, IPlugingName, 0644)
	check(err)
}



func check(e error){
	if e != nil {
		panic(e)
	}
}

/*
* Question and answer
*/
func qAndA(quesion string) string{
	fmt.Println(quesion)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan(){
		return scanner.Text()
	}

	return ""
}