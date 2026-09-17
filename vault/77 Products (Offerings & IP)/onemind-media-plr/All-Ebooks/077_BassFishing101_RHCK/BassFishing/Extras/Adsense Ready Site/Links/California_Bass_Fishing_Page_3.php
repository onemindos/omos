<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title>Links | California Bass Fishing Page 3</title>
		<META http-equiv="Content-Type"	content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="California Bass Fishing Related Links - Page 3 - ">
		<meta name="Keywords"			content="<?=$Keywords;?>">
		<meta name="Distribution"		content="global">
		<meta name="Publisher"			content="<?=$Domain;?>">
		<meta name="Rating"				content="General">
		<meta name="Revisit-after"		content="5 days">
		<meta name="Robots"				content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0">
			<tr>
				<td class="tdHeader" colspan="2">
					<p style="margin-bottom: 5; margin-right: 10;">
						<?php include("../Includes/Menu.php"); ?>
					</p>
				</td>
			</tr>
			<tr>
				<td class="tdContent">
					<br>
					<p>
						<b>California Bass Fishing Related Links - Page 3 </b>
						<br>
						<b>
							<small>
								
							</small>
						</b>
						<br>
						<br>
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<ul class='ulCategory'><li><A href="http://www.probassguides.com/guidelist.php?statename=California%20-%20N" target='_blank'>California - N Bass Fishing Guides</A><br><p>California - N Bass Guides Western Outfitting - 875 W. School St. Cotati, California - N 94931 We specialize in custom guided fishing vacations. Our guides are well versed in many types of fishing ... </P></li></ul>
						<br>
						<br>
						<Center><a href='California_Bass_Fishing_Page_2.php'>Back</a> | <a href='California_Bass_Fishing_Page_1.php'>1</a> | <a href='California_Bass_Fishing_Page_2.php'>2</a> | 3 | <font color='CCCCCC'>Next</font><br><br></center>
					</p>
					<br>
					<?php	
						if ($DisplayAmazon)
							{
								echo "<br><center>";
								include("../Includes/Amazon_468x60.php");
								echo "</center><br>";
							}
					?>					
				</td>
				<td class="tdRight">
					<?php include("../Includes/Navigation.php"); ?>
					<?php include("../Includes/Google_160x600.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdRow" colspan="2">
					<?=$Category;?> News and Events
				</td>
			</tr>
			<tr>
				<td colspan="2">
					<?php include("../Includes/Google_Search.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="2">
					<br>
					<?php include($Domain."Includes/NewsFeed.php?CAT=California+Bass+Fishing"); ?>
					<?php	
						if ($DisplayAmazon)
							{
								echo "<br><center>";
								include("../Includes/Amazon_728x90.php");
								echo "</center><br>";
							}
					?>
				</td>
			</tr>
			<tr>
				<td class="tdFooter" colspan="2">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>

