<?php 

 	header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
	header('Cache-Control: no-store, no-cache, must-revalidate');
	header('Cache-Control: post-check=0, pre-check=0', FALSE);
	header('Pragma: no-cache');

	$UserName				= "admin";
	$Password				= "password";

	$Domain					= "http://".$_SERVER["HTTP_HOST"].rtrim(str_replace($DOCUMENT_ROOT, "", dirname(dirname($PHP_SELF))), "/")."/";
	$SiteName				= "";			
	
	$MainTitle				= "";
	$SubTitle				= "";
	
	$Category				= "Bass Fishing";
	$NewsFeed				= "Bass+Fishing";
			
	$Keywords				= "Basic Guides On Bass Fishing, Bass Fisherman's Techniques, Bass Fishing Big Catch, Bass Fishing Central Florida, Bass Fishing Charter, Bass Fishing Christmas Card, Bass Fishing Facts, Bass Fishing Game Online, Bass Fishing Home Page, Bass Fishing Homepage, Bass Fishing In FLorida Recreation, Bass Fishing In Pueto Vallarta, Bass Fishing Lures, Bass Fishing Reports, Bass Fishing Reports, Bass Fishing Tackle, Bass Fishing Techniques, Bass Fishing Tips, Bass Fishing Tournaments, Bass Fishing Tournaments In Florida, Bass Fishing With Family In FLorida, Best Fishing Areas, Best Guide For Bass Fishing, Best Place To FIsh In Florida, California Bass Fishing, Canada Bass Fishing, Catch Bass Fish Of Lifetime, Effective Bass Fishing Lures, Experiencing Bass Fishing In Florida, Fall Bass Fishing, Florida Bass Fishing, Florida Bass Fishing Center, Free Bass Fishing Games, Georgia Bass Fishing, Hot Bass Fishing Suggestions, Lake Okeechobee Bass Fishing, Linking To Bass Fishing, Mexico Bass Fishing, North Carolina Bass Fishing, Ontario Bass Fishing, Peacock Bass Fishing, Practical Bass Fishing, Pro Bass Fishing, Quest For Smallmouth Bass, Smallmouth Bass, Smallmouth Bass Fishing In Canada, Smallmouth Bass In Canada, Striped Bass Fishing, Texas Bass Fishing, Thousand Islands Bass Fishing";
	
	$GooglePubID			= ""; 
	$GoogleChannel			= ""; 
	$GoogleLinkColor		= "6485B8";
	$GoogleURLColor			= "9EAA86";	

	$AmazonID				= "";
	$AmazonMode				= "Sporting";
	$AmazonSearch			= "Bass Fishing";
	$AmazonLinkColor		= "6485B8";
	$AmazonTextColor		= "9EAA86";
	$DisplayAmazon			= "";

	$CopyRight				= "All images, text and code found on this web site are &copy; ".date("Y").",".$SiteName." - All rights reserved worldwide..."; 

	$Name					= "";
	$Address1				= "";
	$Address2				= "";
	$City					= "";
	$State					= "";
	$ZipCode				= "";
	$Country				= "";
	$Telephone				= "";
	$EmailAddress			= "";

	$ShowGoogle				= "checked";
	$ShowFireFox			= "checked";
	$ShowSearch				= "checked";

	$RotateArticles			= "checked";	

?>