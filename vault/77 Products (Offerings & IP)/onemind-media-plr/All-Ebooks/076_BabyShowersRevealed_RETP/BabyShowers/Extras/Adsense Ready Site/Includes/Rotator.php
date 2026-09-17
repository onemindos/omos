<?php 
				
	$dirList					= opendir(getcwd());  	
	$arrList					= array(); 				
	
	while (($dirPart			= @readdir($dirList)) == true) 
		{  			
			if (substr($dirPart, (strlen($dirPart) -4), 4) == ".php")	
				{
					$arrList[]	= $dirPart; 
				}
		} 	
	
	$rndKey						= rand(0,sizeof($arrList)-1); 	
	
	header("Location: ".$arrList[$rndKey]);
?>
