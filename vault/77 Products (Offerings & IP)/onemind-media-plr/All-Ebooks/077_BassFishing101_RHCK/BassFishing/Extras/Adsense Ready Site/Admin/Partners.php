<?php include("Session.php"); ?>
<?php

	$strPartners 		= "disabled";
	
	if (isset($_POST['Partners']))
		{	
			$Partners	= $_POST['Partners'];
			$fp			= fopen("../Includes/Partners.php", "w");
			$fw			= fwrite( $fp, $Partners );
			fclose($fp);

			$message  = "Your Partner's Webpage has been saved!";
			$strAlert = "<table border=0 cellpadding=0 cellspacing=0 width=100% style='border-color:#000000;border-style:solid;border-width:1;' bgcolor=#FFFFCC><tr><td align=center>$message</td></tr></table>";				
		}
	else
		{
			$Partners = file_get_contents("../Includes/Partners.php");
		
		}
?>
	<?php include("Header.php"); ?>
		<form method="post" action="Partners.php" name="frmData">
			<center>
				<table class="tblForm">
					<tr>
						<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
							Partner's Directory
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<?php 
								if($message)
									{
										echo $strAlert;
									}
								else
									{
							?>
										&nbsp;
							<?php
									}	
							?>
						</td>				
						</td>				
					<tr>
					<tr>
						<td class="tdFieldValue" colspan="2">
							<textarea name="Partners" rows="15" style="width: 100%;"><?=stripslashes($Partners);?></textarea>
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<hr size="0">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
						</td>				
						<td class="tdFieldValue">
							<input class="iButton" onClick="this.disabled='true'; frmData.submit();" type="submit" value=" Save Partner's Directory ">
						</td>				
					<tr>
				</table>
			</center>
	<?php include("Footer.php"); ?>