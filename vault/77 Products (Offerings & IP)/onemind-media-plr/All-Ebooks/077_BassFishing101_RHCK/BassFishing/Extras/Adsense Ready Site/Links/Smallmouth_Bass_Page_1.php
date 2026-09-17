<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title>Links | Smallmouth Bass Page 1</title>
		<META http-equiv="Content-Type"	content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="Smallmouth Bass Related Links - Page 1 - ">
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
						<b>Smallmouth Bass Related Links - Page 1 </b>
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
						<ul class='ulCategory'><li><A href="http://www.smallmouth.ca/" target='_blank'>Smallmouth Bass</A><br><p>... site regarding the Smallmouth Bass fishing Canadian style. All kinds of information available from tournaments to product reviews. Visit us today! A place to acquire more knowledge of the great fishing ... </P></li><li><A href="http://www.smallmouthquest.com/" target='_blank'>Smallmouth bass fishing with guide Darl Black</A><br><p>Fishing tips, resorts, guides and resources from bass fishing guide and outdoor writer Darl Black.</P></li></ul>
						<br>
						<br>
						<Center><font color='CCCCCC'>Back</font> | 1 | <a href='Smallmouth_Bass_Page_2.php'>2</a> | <a href='Smallmouth_Bass_Page_3.php'>3</a> | <a href='Smallmouth_Bass_Page_2.php'>Next</a><br><br></center>
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
					<?php include($Domain."Includes/NewsFeed.php?CAT=Smallmouth+Bass"); ?>
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

