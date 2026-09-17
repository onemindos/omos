<?php include("Header.php"); ?>
	<?php		
	
		error_reporting(0);
		
		chmod(realpath("../Includes/Configure.php"), 0666);
		chmod(realpath("../Sitemap.xml"), 0666);
						
		if(!is_writable("../Includes/Configure.php"))
			{
				$ERR_Configure = true;
			}
		if(!is_writable("../Sitemap.xml"))
			{
				$ERR_Sitemap = true;
			}
		if($ERR_Configure or $ERR_Sitemap)
			{
	?>
				<center>
					<table class="tblForm">
						<tr>
							<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_005.jpg); background-repeat: repeat-x; height: 14px;">
								&nbsp; 
							</td>				
						</tr>
						<tr>
							<td colspan="2">
								<br>
								<img align="left" src="http://www.niche-mania.com/Images/NM_004.jpg">
								<p>
									<b>
										W A R N I N G !
									</b>
									<br>
									<br>
									In order to use this program, you must first set the file permissions on the following files to <b>666</b>, which will allow it to both read and write to the files:
									<br>
								</p>
								<ul>
									<?php
										if($ERR_Configure)
											{
									?>
												<li>
													<?php echo $_SERVER["HTTP_HOST"]; ?>/Includes/<b style="color: #C80000;">Configure.php</b>
												</li>
									<?php	
											}
										if ($ERR_Sitemap)	
											{
									?>
												<li>
													<?php echo $_SERVER["HTTP_HOST"]; ?>/<b style="color: #C80000;">Sitemap.xml</b>
												</li>
									<?php	
											}
									?>
								</ul>
								<p>
									You can most easily do this using an FTP program such as:
								</p>
								<ul>
									<li>
										<b>
											Cute FTP 
										</b>
										 (<a href="http://www.cuteftp.com" target="_new">http://www.cuteftp.com</a>)
									</li>
									<li>
										<b>
											Filezilla
										</b>
										 (<a href="http://filezilla.sourceforge.net/" target="_new">http://filezilla.sourceforge.net</a>)
									</li>
									<li>
										<b>
											WS FTP
										</b>
										 (<a href="http://www.ipswitch.com" target="_new">http://www.ipswitch.com</a>)
									</li>
								</ul>
								<p>
									After doing so, please return to this page to login!									
								</p>
								<br>
							</td>				
						</tr>
						<tr>
							<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_005.jpg); background-repeat: repeat-x; height: 14px;">
								&nbsp; 
							</td>				
						</tr>
					</table>
				</center>
	<?php
			}
		else
			{
	?>				
				<form method="post" action="Login.php" name="frmData">
					<center>
						<table class="tblForm">
							<tr>
								<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
									CONTROL PANEL LOGIN&nbsp; 
								</td>				
							<tr>
							<tr>
								<td colspan="2">
									<small>
										&nbsp;
									</small>
								</td>				
							<tr>
							<tr>
								<td class="tdFieldName">
									User Name
								</td>				
								<td class="tdFieldValue">
									<input class="iText" style="width: 100%;" type="text" name="UserName">
								</td>				
							<tr>
							<tr>
								<td class="tdFieldName">
									Password
								</td>				
								<td class="tdFieldValue">
									<input class="iText" style="width: 100%;" type="password" name="Password">
								</td>				
							<tr>
							<tr>
								<td colspan="2">
									<hr noshade size="1">
								</td>				
							<tr>
							<tr>
								<td class="tdFieldName">
									<p class="pError">
										<?php if (isset($_GET['Status'])) {echo " LOGIN ERROR! ";} ?>
									</p>
								</td>				
								<td class="tdFieldValue">
									<input class="iButton" onClick="this.disabled='true'; frmData.submit();" type="submit" value=" Login ">	
								</td>				
							<tr>
						</table>
					</center>
				</form>
	<?php
			}
	?>
<?php include("Footer.php"); ?>