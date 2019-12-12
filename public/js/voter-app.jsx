
function toggle(Inner) {
  return props => {
    const [toggled, setToggled] = React.useState(false);
    const onClick = React.useCallback(() => setToggled(t => !t), []);
    return <Inner {...props} toggled={toggled} onClick={onClick} />;
  };
}

const TensionDemo = toggle(({toggled, ...props}) => (
  <section>
    <h2>Translate with high tension</h2>
    <Animate translateX={toggled ? 200 : 0} tension={500}>
      <button className="c4" {...props}>
        Click Me
      </button>
    </Animate>
  </section>
));

const App = () => (
  <div>
    <h1>react-rebound demos</h1>
    <p>
      See <a href="https://github.com/steadicat/react-rebound">react-rebound on GitHub</a> for code
      and instructions.
    </p>
    <p>
      Source for these examples is{' '}
      <a href="https://github.com/steadicat/react-rebound/blob/master/demo.tsx">here</a>.
    </p>
    
    <TensionDemo />
    
  </div>
);

ReactDOM.render(<App />, document.getElementById('app'));