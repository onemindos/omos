# How to Access Your Hermes Agent Anywhere

> **Title:** How to Access Your Hermes Agent Anywhere
> **Author:** tonbi ([@tonbistudio](https://x.com/tonbistudio))
> **Source:** [x.com/tonbistudio/status/2060945344092148007](https://x.com/tonbistudio/status/2060945344092148007?s=12)
> **Type:** X Article | **Posted:** 2026-05-31 | **Ingested:** 2026-06-07

---

Your local AI agent is probably trapped

When people first start using AI agents, they usually rush toward the fun stuff.

Skills. Plugins. Cron jobs. Tool calls. Weird little automations that run in the background. That stuff is great. It's the fun part.

But there is a less exciting layer that matters just as much: the environment around the agent.

Because for a lot of people, their local agent is trapped on one machine.

Maybe it lives on your desktop because that is the strongest computer you own. Maybe it runs on a Mac Mini in the corner. Maybe it's on an old laptop you repurposed into a little always-on agent box.

That works fine while you are sitting in front of it.

But what happens when you leave?

What happens when you start a long research task, close your laptop, lose SSH, switch Wi-Fi networks, or want to check progress from your phone?

Suddenly your “autonomous agent” is not very autonomous.

It is a process running on one computer that you can only comfortably access from one chair.

That is the problem this stack solves.

The goal is simple:

I want to access my Hermes agent from anywhere.

I do not want to expose ports to the public internet.

I do not want fragile SSH sessions that die when my connection drops.

I do not want to restart a task every time I move from my desktop to my laptop to my phone.

The stack I have been using is:

Tailscale for reach.

Termius for control.

tmux for persistence.

Hermes is the agent. These three tools are the infrastructure layer around it.
Now I'm gonna tell you about it, but if you prefer to watch, check out this video:
https://youtu.be/KASkBTo3hI4

Tailscale is what makes the whole thing possible.

It creates a private mesh network across your devices using WireGuard. Your desktop, laptop, phone, server, WSL environment, whatever you are using, all get private addresses and readable names.

So instead of exposing your machine to the internet, opening ports, setting up port forwarding, or messing with random tunnel hacks, your devices can talk to each other privately.

That matters a lot when the thing you are accessing is an AI agent with shell access.

Port forwarding a dev server for ten minutes is one thing.

Leaving an agent reachable on the public internet is another. That is how you end up in “future cautionary blog post” territory.

With Tailscale, there is no public IP you need to expose. No router config. No open SSH port sitting on the internet. Your devices connect through your private Tailnet.

If Hermes is running inside WSL, install Tailscale inside WSL. Not just on the Windows host. Install it where the agent actually runs.

The basic setup is pretty painless:
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --hostname=hermes-pc
tailscale status
Then install Tailscale on your phone and laptop, log into the same account, and now those devices can reach the Hermes machine by name.

That solves the reach problem.

Next is control.

For that, you need SSH.

SSH is just the way you log into another machine’s command line over an encrypted connection. The Hermes machine runs the server. Your laptop or phone is the client.

On the machine running Hermes:
sudo apt install openssh-server
sudo systemctl enable --now ssh
sudo systemctl status ssh
Now you can SSH into that machine over Tailscale.

On a laptop, this is easy. Use your terminal.

On a phone, most SSH apps are awful.

This is where Termius comes in.

Termius is a mobile SSH client that does not make me want to throw my phone into the ocean. It gives you a usable keyboard bar, SSH keys, host management, port forwarding, a file browser, and enough quality-of-life stuff that controlling a remote shell from your phone is actually tolerable.

You create a host in Termius using the Tailscale hostname or private IP, set up SSH key auth, and connect.

At that point, your phone can open a shell on the machine running Hermes.

You can run:
hermes status
hermes chat
You can also set up local port forwarding so the Hermes dashboard is available from your phone browser.

That is the moment the whole setup starts to feel different.

Your local agent is no longer just sitting on your desktop. It is in your pocket.

But there is still one missing piece.

SSH gives you access, but a normal SSH session is fragile.

If you connect over SSH, start Hermes, then lose connection, close your laptop, switch networks, or walk into bad Wi-Fi, that session can die.

That is annoying for normal terminal work.

It is much worse for agent work, because agents often do long running tasks.

Research something. Watch a process. Run a job. Wait for feedback. Continue a conversation.

You do not want all of that tied to whether your phone kept a clean connection.

That is why you use tmux.

tmux owns the shell.

SSH becomes just a window into it.

So you start a tmux session on the Hermes machine:
sudo apt install tmux
tmux new -s hermes
hermes chat
Now Hermes is running inside tmux.

If your SSH connection drops, tmux keeps running.

If you close your laptop, tmux keeps running.

If your phone loses signal, tmux keeps running.
As long as the main device running the agent is running, your agent will keep working regardless.

Later, from any device on your Tailnet, you SSH back in and reattach:
tmux attach -t hermes
Same session. Same agent. Same work.

That is the important part.

You can start a long task on your desktop, leave the house, open Termius from your phone, reattach to the tmux session, see what Hermes has done, give it another instruction, then later continue from your laptop.

One agent. One persistent workspace. Multiple devices.

No public ports.

No restarting every time your connection changes.

No babysitting a fragile terminal session.

This is the part I think people underestimate with agents.

The model matters. The tools matter. The prompts matter.

But the boring infrastructure around the agent matters too.

If an agent is supposed to work for a while, you need a way to leave and come back.

If it is supposed to be useful throughout the day, you need access from more than one device.

If it has shell access, you probably should not expose it to the public internet like a maniac.

Tailscale, Termius, and tmux are not the only way to do this.

You could use Cloudflare Tunnel. You could use a VPS with reverse SSH. You could build something more custom. There are plenty of options.

But this stack is beginner friendly, fast to set up, and good enough for most people running a local agent.

The mental model is clean:

Tailscale solves reach.

Termius solves control.

tmux solves persistence.

Hermes does the work.

One honest caveat: this stack secures how you reach the agent. It does not sandbox what the agent can do after you reach it.

That distinction matters.

Tailscale can keep you from exposing SSH to the internet. SSH keys can gate access. tmux can keep sessions alive.

But if your agent has access to your files, API keys, shell, browser, or private repos, you still need to think carefully about permissions.

This setup does not magically solve prompt injection.

It does not turn a dangerous agent into a safe one.

It just gives you a much better way to access the machine where that agent runs.

If you want to harden it further, look into Tailscale ACLs. By default, devices on your Tailnet can usually talk to each other. ACLs let you decide which devices can reach which machines and services.

That is probably overkill for some personal setups, but worth understanding if you are putting more serious workflows behind this.

The bigger point is this:

Local agents become much more useful when they stop being tied to one screen.

I want to start a task on my desktop, check it from my phone, continue from my laptop, and not care where the process originally started.

I want the agent to keep working when I disconnect.

I want the dashboard in my pocket.

I want zero public internet exposure.

That is what this stack gives you.
The fun part of agents is skills, plugins, tools, and weird automations.

The useful part is when the agent is always reachable, still running, and not hanging off the open internet by a thread.

Boring infrastructure is what makes the fun stuff actually usable.
