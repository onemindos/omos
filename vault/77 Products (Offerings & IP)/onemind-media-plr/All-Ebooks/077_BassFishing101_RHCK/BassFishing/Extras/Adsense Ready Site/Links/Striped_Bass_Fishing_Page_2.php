<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title>Links | Striped Bass Fishing Page 2</title>
		<META http-equiv="Content-Type"	content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="Striped Bass Fishing Related Links - Page 2 - ">
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
						<b>Striped Bass Fishing Related Links - Page 2 </b>
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
						<ul class='ulCategory'><li><A href="http://www.striped-bass-fishing.com/bboard/viewforum.php?f=40" target='_blank'>Striped-Bass-Fishing.com - South Carolina </A><br><p>striped-bass-fishing.com - Striper tournament fisherman [ Goto page: 1 , 2 ] 24 panfish16 690 Wed Aug 03, 2005 4:00 pm skeet LAKE JOCASSEE: Trout: Good Try drifting large minnows from 0 NewsReader 98 Sat Jul 30 ... </P></li></ul>
						<br>
						<br>
						<Center><a href='Striped_Bass_Fishing_Page_1.php'>Back</a> | <a href='Striped_Bass_Fishing_Page_1.php'>1</a> | 2 | <a href='Striped_Bass_Fishing_Page_3.php'>3</a> | <a href='Striped_Bass_Fishing_Page_3.php'>Next</a><br><br></center>
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
					<?php include($Domain."Includes/NewsFeed.php?CAT=Striped+Bass+Fishing"); ?>
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

