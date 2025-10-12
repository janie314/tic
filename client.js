while (true) {
  const res = await fetch("http://192.168.1.82/tic/games").then((r) =>
    r.json()
  );

  for (const game of res) {
    console.log(game);
    if (game.turn === "jane") {
      for (let i = 0; i < 3; i++) {
        for (let j = 0; j < 3; j++) {
          if (game.state[i][j] === "_") {
            fetch("http://192.168.1.82/tic/move", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({
                id: game.id,
                user: "jane",
                row: i,
                col: j,
              }),
            });
          }
        }
      }
    }
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
}
