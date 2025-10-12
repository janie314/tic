const [_t, _s, host, user] = Bun.argv;

while (true) {
    let res;
    const res1 = await fetch(`http://${host}/tic/games`);
    if (res1.ok) {
        res = await res1.json()
    } else {
        continue;
    }
    console.log("loaded games");
    for (const game of res) {
        console.log('game', game.id, 'turn', game.turn);
        if (game.turn === user) {
            while (true) {
                const i = Math.floor(Math.random() * 3);
                const j = Math.floor(Math.random() * 3);
                if (game.state[i][j] === "_") {
                    const body = {
                        id: game.id,
                        user: user,
                        row: i,
                        col: j,
                    };
                    console.log(game);
                    console.log(body);
                    fetch(`http://${host}/tic/move`, {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(body),
                    });
                    break;
                }
            }
        }
    }
    await new Promise(resolve => setTimeout(resolve, 100));
}
