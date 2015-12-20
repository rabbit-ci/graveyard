import DS from 'ember-data';

export default DS.Model.extend({
  stdio: DS.attr('string') // TODO: this needs to be parsed with
  // ansi_up since it will include ANSI escape codes
});
