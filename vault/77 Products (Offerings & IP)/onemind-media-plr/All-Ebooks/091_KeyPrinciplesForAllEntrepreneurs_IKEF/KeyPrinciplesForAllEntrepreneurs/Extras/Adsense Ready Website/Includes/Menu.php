<?php

	// -----------------------------------------------------------------
 	// You may edit this file to add additional links to other areas of
 	// your site. This file is included on every page by default. 
 	// © 2000 - 2006, Niche Mania - All Rights Reserved Worldwide...
 	// -----------------------------------------------------------------
	
	error_reporting(0);	
	
	if ($Page == "Articles")
		{
			echo "<b>".$Category." Articles</b> | ";
		}
	else
		{
			echo "<a href='../Articles/index.php'>".$Category." Articles</a> | ";
		}

	if ($Page == "Links")
		{
			echo "<b>".$Category." Links</b> | ";
		}
	else
		{
			echo "<a href='../Links/index.php'>".$Category." Links</a> | ";
		}

	if ($Page == "Partners")
		{
			echo "<b>".$Category." Partners</b> | ";
		}
	else
		{
			echo "<a href='../Partners/index.php'>".$Category." Partners</a> | ";
		}

	if ($Page == "Sitemap")
		{
			echo "<b>".$Category." Sitemap</b>";
		}
	else
		{
			echo "<a href='../Sitemaps/index.php'>".$Category." Sitemap</a>";
		}
?>