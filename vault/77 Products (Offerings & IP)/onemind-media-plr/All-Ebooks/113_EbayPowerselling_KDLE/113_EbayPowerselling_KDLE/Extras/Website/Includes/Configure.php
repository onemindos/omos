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

	$Category				= "Ebay";
	$NewsFeed				= "Ebay";
	
	$Keywords				= "19 Free Ebay Secrets, Becoming A Ebay Powerseller, Best Products To Sell On Ebay , Best Things To Sell On Ebay, Buy To Sell On Ebay , Ebay Business Opportunity, Ebay For Newbies, Ebay For Starters, Ebay Marketplace Selling , Ebay Opportunities In Franchising, Ebay Secrets, Ebay Selling Secrets Unveiled, Ebay Selling Software, Ebay Selling Tips, Ebay Selling Tips For Everyone, Ebay Selling Tips From Power Buyer, Ebay Selling Tools, Ebay Training, Ebay Wholesale sellers, Free Ebay Business Training, Free Ebay Selling Tips, Guides To Selling On Ebay, Hot Selling Items On Ebay, How To Easily Sell On Ebay, How To Make Money On Ebay, How To Start An Ebay Business, Learn How To Sell On Ebay, Lucrative Ebay Opportunities, Make Money On Ebay Without Dropshipping, Make Money With An Ebay Business , Making Extra Money On Ebay, Making Money On Ebay, Making Money With An Ebay Business, Perfect Way To Start An Ebay Business, Products To Sell On Ebay, Secrets To A Successful Ebay Business, Sell On Ebay The Right Way , Selling cds And Dvds On Ebay, Selling On Ebay After The Holidays,  Start Making Money With Ebay, Start Your Own Ebay Business, Steps To Ebay Selling Success, Think Like A Ebay Powerseller, Top Ebay Seller Secrets, Top Ebay Selling Items For 2005, Ways To Make Money On Ebay, What It Takes To Become A Ebay Powerseller, What Sells On Ebay, What To Sell On Ebay, Wholesale Selling On Ebay";
	
	$GooglePubID			= "";
	$GoogleChannel			= "";
	$GoogleLinkColor		= "05019C";
	$GoogleURLColor			= "98CD00";

	$AmazonID				= "";
	$AmazonMode				= "Books";
	$AmazonSearch			= "Ebay";
	$AmazonLinkColor		= "05019C";
	$AmazonTextColor		= "98CD00";
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