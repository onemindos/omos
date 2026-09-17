<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title>Links | Ebay For Newbies Page 1</title>
		<META http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Creator" content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright" content="<?=$CopyRight;?>">
		<meta name="Description" content="Ebay For Newbies Related Links - Page 1 - ">
		<meta name="Keywords" content="<?=$Keywords;?>">
		<meta name="Distribution" content="global">
		<meta name="Publisher" content="<?=$Domain;?>">
		<meta name="Rating" content="General">
		<meta name="Revisit-after" content="5 days">
		<meta name="Robots" content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0">
			<tr>
				<td class="tdHeader" colspan="3">
					<p style="margin-bottom: 5; margin-left: 10;">
						<?php include("../Includes/Menu.php"); ?>
					</p>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include("../Includes/Google_728x15_TextLinks.php"); ?>
					<hr>
				</td>
			</tr>
			<tr>
				<td class="tdleft">
					<br>
					<?php include("../Includes/Navigation.php"); ?>
				</td>
				<td class="tdContent">
					<br>
					<p>
						<b>Ebay For Newbies Related Links - Page 1 </b>
						<br>
						<b><small> </small></b>
						<br>
						<br>
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<ul class='ulCategory'><li><A href="http://ebaypink.blogspot.com/2005/11/long-weekends-and-newbies.html" target='_blank'>eBay pink: Long Weekends and Newbies</A><br><p>One thing I have noticed this weekend is the number of eBay Newbies having a go. They seem ot have been brought out by the extra time on pepole's hands. It's good to see. And it's nice when they bid ... </P></li><li><A href="http://www.freebiddingtools.com/help" target='_blank'>A Beginner's Guide to eBay - eBay for Newbies</A><br><p>... How to use eBay - A Beginner's Guide to eBay New to eBay How does eBay work? How to Register Buying Overview Buying with Confidence How to Buy How to Bid and Buy Successful Searching How to Pay for your Item ... </P></li></ul>
						<br>
						<br>
						<Center><font color='CCCCCC'>Back</font> | 1 | <a href='Ebay_For_Newbies_Page_2.php'>2</a> | <a href='Ebay_For_Newbies_Page_3.php'>3</a> | <a href='Ebay_For_Newbies_Page_2.php'>Next</a><br><br></center>
					</p>
					<br>
					<?php if ($DisplayAmazon) { echo "<br><center>"; 
					include("../Includes/Amazon_468x60.php"); echo "</center><br>"; } ?>
				</td>
				<td class="tdRight">
					<?php include("../Includes/Google_160x600.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="3" class="tdRow">
					<?=$Category;?> News and Current Events
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include("../Includes/Google_Search.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include($Domain."Includes/NewsFeed.php?CAT=$NewsFeed"); ?> <?php if 
					($DisplayAmazon) { echo "<br><center>"; 
					include("../Includes/Amazon_728x90.php"); echo "</center><br>"; } ?>
				</td>
			</tr>
			<tr>
				<td class="tdFooter" colspan="3">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All 
					Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal 
						Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>

