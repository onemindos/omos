<?php include("Session.php"); ?>
<?php

	$strFileName			= "Sitemap.xml";
	$strDomain				= "http://".$_SERVER["HTTP_HOST"];
            
	$strXMLHeader			= "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n";
	$strXMLNameSpace		= "  <urlset xmlns=\"http://www.google.com/schemas/sitemap/0.84\">\n";
	$strXMLCloseNS			= "</urlset>";

	$Articles				= CreateNodeString("../Articles/", "Articles/");
	$Legal					= CreateNodeString("../Legal/", "Legal/");
	$Links					= CreateNodeString("../Links/", "Links/");
	$Partners				= CreateNodeString("../Partners/", "Partners/");
	$Sitemaps				= CreateNodeString("../Sitemaps/", "Sitemaps/");

	$strXMLDoc				= $strXMLHeader . $strXMLNameSpace . $Articles . $Legal . $Links . $Partners . $Sitemaps . $strXMLCloseNS;

	WriteXMLDocument($strXMLDoc);
	
	$_SESSION['XMLSitemap']	= true;
		
	header('location: Sitemap.php');
		            
	function CreateNodeString($path, $RealPath)
		{
			global $strDomain;
            
			$num			= array("0.5","0.6","0.7","0.8","0.9","1.0");
			$strNode		= "";
			$dtString		= get_iso_8601_date(time());
			$Dir			= opendir($path);
            
			while ($file = readdir($Dir))
				{                
					if($file != "." and $file != "..")
						{					
							$ext	= explode(".",$file);
							$url	= $_SESSION['Domain'].$RealPath.$file;
	                
							shuffle($num);
							$strNode = $strNode . "<url>\n<loc>$url</loc>\n<lastmod>$dtString</lastmod>\n<changefreq>weekly</changefreq>\n<priority>$num[0]</priority>\n</url>\n";
						}
				}
            
			closedir($Dir);
            
			return $strNode;
		}
		
	function WriteXMLDocument($strXML)
		{
			global $strFileName;
            
			$fp				= fopen("../".$strFileName, "w" );
			fwrite($fp, "");
			fclose( $fp );

			$fp				= fopen("../".$strFileName, "a");
			$fw				= fwrite( $fp, "$strXML" );
			fclose( $fp );
		}
        
	function get_iso_8601_date($int_date) 
		{
			$date_mod		= date('Y-m-d', $int_date)."T".date('H:i:s', $int_date);
			$pre_timezone	= date('O', $int_date);
			$time_zone		= substr($pre_timezone, 0, 3).":".substr($pre_timezone, 3, 2);
			$date_mod	   .= $time_zone;
			
			return $date_mod;
        }        
?>