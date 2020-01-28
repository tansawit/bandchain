[@bs.val] external document: Js.t({..}) = "document";

let style = document##createElement("style");
document##head##appendChild(style);
style##innerHTML #= AppStyle.style;

Axios.setRpcUrl("https://d3n.bandprotocol.com/");

TimeAgos.setMomentRelativeTimeThreshold();

ReactDOMRe.render(<GlobalContext> <App /> </GlobalContext>, document##getElementById("root"));
